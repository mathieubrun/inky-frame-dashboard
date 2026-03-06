import gc
import machine
import network
import time
import inky_frame
import pngdec
import urequests
from picographics import PicoGraphics, DISPLAY_INKY_FRAME_7_SPECTRA as DISPLAY_TYPE

# Import local configuration
try:
    import env
except ImportError:
    print("Error: env.py not found. Please create one based on .env.template.py")
    inky_frame.sleep_for(30) # Default sleep to prevent battery drain

# Hardware Pins for VSYS (Battery Monitoring)
PIN_VSYS_ADC = 29
PIN_VSYS_HOLD = 2
PIN_VSYS_READ_EN = 25

# Constants
VOLTAGE_CONVERSION_FACTOR = 3 * 3.3 / 65535
TEMP_IMAGE_FILE = "dashboard.png"

# Initialize Display
display = PicoGraphics(display=DISPLAY_TYPE)
display.set_font("bitmap8")

def get_battery_voltage():
    """Reads the VSYS voltage via ADC to estimate battery level."""
    # Keep the board powered on (essential for battery use)
    hold_vsys_en_pin = machine.Pin(PIN_VSYS_HOLD, machine.Pin.OUT)
    hold_vsys_en_pin.value(True)

    # Prepare ADC for VSYS (Pico W specific)
    vsys_read_en = machine.Pin(PIN_VSYS_READ_EN, machine.Pin.OUT)
    vsys_read_en.value(True)

    vsys_adc = machine.ADC(PIN_VSYS_ADC)

    # Take a few readings and average them for stability
    reading = 0
    for _ in range(10):
        reading += vsys_adc.read_u16()
    reading /= 10
    
    # Restore pins
    vsys_read_en.value(False)

    return reading * VOLTAGE_CONVERSION_FACTOR

def set_led(color, state):
    """Controls the onboard activity LED."""
    # Inky Frame has a single activity LED
    inky_frame.button_a.led_on() if state else inky_frame.button_a.led_off()

def connect_wifi(ssid, password, timeout=15):
    """Connects to the specified Wi-Fi network with a timeout."""
    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)
    
    if wlan.isconnected():
        return True

    print(f"Connecting to WiFi: {ssid}")
    wlan.connect(ssid, password)

    start_time = time.time()
    # Blink LED while connecting
    led_state = False
    while not wlan.isconnected():
        if time.time() - start_time > timeout:
            print("WiFi Connection Timeout")
            wlan.active(False)
            set_led("network", False)
            return False
        set_led("network", led_state)
        led_state = not led_state
        time.sleep(0.5)

    set_led("network", True) # Solid on connected
    print(f"Connected. IP: {wlan.ifconfig()[0]}")
    return True

def sleep():
    """Puts the Inky Frame into deep sleep."""
    set_led("all", False) # Ensure all LEDs are off before sleep
    try:
        sleep_mins = env.SLEEP_MINUTES
    except AttributeError:
        sleep_mins = 30 # Default if not in env
    
    print(f"Going to sleep for {sleep_mins} minutes...")
    inky_frame.sleep_for(sleep_mins)

def fetch_image(url, filename):
    """Fetches the image from the specified URL and saves it to a file."""
    print(f"Fetching image from {url}")
    try:
        # Important: Setting timeout is crucial for battery devices
        response = urequests.get(url, timeout=30)
        
        if response.status_code == 200:
            with open(filename, 'wb') as f:
                f.write(response.content)
            response.close()
            print("Image downloaded successfully.")
            return True
        else:
            print(f"HTTP Error: {response.status_code}")
            response.close()
            return False
            
    except Exception as e:
        print(f"Fetch failed: {e}")
        return False

def render_image(filename, voltage):
    """Decodes the PNG file and updates the e-ink display."""
    print(f"Rendering image: {filename}")
    try:
        # Create PNG decoder
        png = pngdec.PNG(display)
        
        # Open and decode the image to the display buffer
        png.open_file(filename)
        png.decode(0, 0)
        
        # Check battery threshold and overlay if necessary
        try:
            threshold = env.BATTERY_THRESHOLD
        except AttributeError:
            threshold = 3.4 # Default
            
        if voltage < threshold:
            print("Drawing low battery warning overlay...")
            display.set_pen(2) # Red (assuming Spectra 6 palette mapping: 0=Black, 1=White, 2=Red)
            display.rectangle(650, 10, 140, 40)
            display.set_pen(1) # White
            display.set_font("bitmap8")
            display.text("LOW BATT", 660, 20, scale=3)
        
        # Return true here, the actual display.update() will be called in the main loop
        return True
    except Exception as e:
        print(f"Render failed: {e}")
        return False

def draw_error_screen(message):
    """Draws an error message to the display."""
    display.set_pen(1) # White
    display.clear()
    
    display.set_pen(0) # Black
    display.set_font("bitmap8")
    
    display.text("Update Failed", 10, 10, scale=4)
    display.text(message, 10, 50, wordwrap=780, scale=2)
    
    # Show time of failure
    year, month, day, dow, hour, minute, second, _ = machine.RTC().datetime()
    time_str = f"Failed at: {year}-{month:02d}-{day:02d} {hour:02d}:{minute:02d}"
    display.text(time_str, 10, 450, scale=2)
    
    display.update()

def main():
    """Main execution loop for the Inky Frame client."""
    print("Inky Frame Dashboard waking up...")
    
    try:
        # 1. Check Battery
        voltage = get_battery_voltage()
        print(f"Battery Voltage: {voltage:.2f}V")
        
        # 2. Connect to Wi-Fi
        if not connect_wifi(env.WIFI_SSID, env.WIFI_PASSWORD):
            print("Failed to connect to WiFi. Sleeping.")
            draw_error_screen("Failed to connect to Wi-Fi network.")
            sleep()
            return

        # 3. Fetch Image
        gc.collect() # Free up memory before fetch
        if not fetch_image(env.DASHBOARD_URL, TEMP_IMAGE_FILE):
            print("Failed to fetch image. Sleeping.")
            draw_error_screen(f"Failed to download dashboard from:\n{env.DASHBOARD_URL}")
            sleep()
            return
            
        # 4. Render Image
        gc.collect() # Free up memory before rendering
        if render_image(TEMP_IMAGE_FILE, voltage):
            print("Updating e-ink display (this will take ~40 seconds)...")
            display.update()
            print("Update complete.")
        else:
            draw_error_screen("Failed to decode and render the downloaded image.")
        
        gc.collect() # Final cleanup
        
    except Exception as e:
        print(f"Unhandled exception: {e}")
        try:
            draw_error_screen(f"System Error: {e}")
        except:
            pass # Failsafe
            
    # 5. Go back to sleep
    sleep()

if __name__ == "__main__":
    main()
