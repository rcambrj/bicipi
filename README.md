# bicipi

status: incomplete

**bicipi**: pronounced _bee-see-pie_, is the fusion of the Spanish word for bike + Raspberry Pi

Connects to a USB Tacx Fortius T1941 trainer and exposes it via Bluetooth Low Energy for use with Zwift, TrainerRoad or MyWhoosh.

## Objectives

* Lightweight and fast:
	* present the device on BLE just like a modern trainer would (no GUI)
* Easy to install:
	* have a downloadable + flashable image for a range of Raspberry Pis
* Support trainers I have access to:
	* Tacx Fortius T1941
* Support popular apps:
	* Zwift, TrainerRoad, MyWhoosh

## Raspberry Pi installation and usage

Supports:

* Raspberry Pi 3
* (probably) Raspberry Pi Zero 2W

> NixOS officially supports only 64-bit ARM, so the downloadable images don't support 32-bit boards such as Raspberry Pi 0, 0W, 1 and 2. You can, however, download the source and build it yourself.

Download the zip for your pi from the releases page and use [etcher](https://etcher.balena.io/) to burn it to an SD card or USB stick. Insert the SD/USB into your pi and give it power. Connect the Tacx via USB. After a minute or so you should be able to kick the pedal to begin the calibration process.

The calibration process by default warms the motor and tyre up for 5 minutes. If the motor is still spinning after 10 minutes, something is wrong.

To find out what's going on, connect your PC directly to the Raspberry Pi via a network cable, configure your PC with address `192.168.24.123`, subnet `255.255.255.0`, and run `ssh bicipi@192.168.24.24`, when prompted for a password, write `password`. You can now scan the logs for the `bicipi` systemd service with `journalctl -xfeu bicipi` and restart it with `sudo systemctl restart bicipi`.

> [!IMPORTANT]
> DO NOT connect this Raspberry Pi to your regular network or the Internet. The password is easily guessed.

If all goes well, after calibration you should see a Bluetooth Low Energy device being advertised called `bicipi`. You can connect to this with TrainerRoad, MyWhoosh or Zwift.

To power the device off, simply remove power. All disks are mounted read only so as to avoid any corruption.

Known issues:

* there is no way to configure any settings right now (do you want this? star the repo to show interest)

## CLI installation and usage

`nix` and `direnv` are recommended in order to use the devshell. You can continue without either of these but you must install the dependencies yourself. Take a look at devshell.nix to infer the necessary dependencies. It only runs on Linux, on MacOS it spits out some errors.

The tool is a simple Go executable, you can run it with: `go run .`.

On a PC with nix, you could consume the package via flake like so:

```
# flake.nix
inputs.bicipi.url = "github:rcambrj/bicipi";

# configuration.nix
environment.systemPackages = [ inputs.bicipi.packages.${pkgs.stdenv.hostPlatform.system}.bicipi ];
```

Once installed or running, `bicipi` presents the following options:

```
Usage of bicipi:
  -bluetooth-name string
    	The bluetooth device name to advertise (default "bicipi")
  -calibrate
    	Whether to enable initial calibration. (--calibrate=false to disable) (default true)
  -calibration-max int
    	How long in seconds before calibration is abandoned. (default 600)
  -calibration-min int
    	How long in seconds to warm up the motor and tyre during calibration. (default 300)
  -calibration-speed int
    	How fast in km/h to spin the tyre during calibration. (default 20)
  -calibration-tolerance int
    	How fussy to be when considering calibration complete. Lower is more fussy. (default 10)
  -loglevel string
    	The log level. May be one of [trace debug info warn error]. (default "info")
  -serial string
    	The serial device to which Tacx motorbrake is connected. (default is to use USB)
  -slow
    	Whether to poll slowly so that logs are easier to follow.
  -weight uint
    	The approximate weight of the rider + bicycle, used only in simulator mode (Zwift / MyWhoosh). (default 80)
```

## Direct to motor serial connection

See [totalreverse's diagram](https://github.com/totalreverse/ttyT1941/wiki#inside-the-t1902-and-t1942) of the `Brake backside socket` for how to connect a serial adapter to the motorbrake.

Remember to put [a pull-up resistor](https://github.com/totalreverse/ttyT1941/issues/7#issuecomment-619587334) on the RX line (motorbrake TX).

On early Raspberry Pis (3 and 02W), direct-to-motorbrake serial connection works only by using a separate USB to serial adapter. I was unable to get a direct-to-motorbrake serial connection working on my Raspberry Pi 3 with the mini-UART (pins 8(14) + 10(15)), and the `ttyAMA0` UART is needed to use Bluetooth.

## Difference between direct to motorbrake vs via headunit

The direct-to-motor serial connection is faster because the headunit introduces a delay in communicating with the motorbrake.

That is, when a command is sent to the headunit, the headunit forwards the command to the motorbrake, and the headunit does not wait for the response from the motorbrake but rather responds immediately with whatever cached response it has already. The next command might have an updated response, or not! Responses appear to be a few hundred milliseconds behind, and cached for somewhere between 1-2 seconds. Despite there being at least two commands (version and control), there appears to be only one cache entry for responses, so the control command might receive a version response, and vice versa.

All this is handled in `bicipi`, but if you want the fastest updates from your trainer, use direct to motorbrake serial rather than USB.

None of the extra features provided by the headunit are implemented in `bicipi`, eg. steering, buttons.

## FortiusAnt

FortiusAnt by WouterJD is a comprehensive software package. It connects to a wide range of hardware, provides a GUI, exposes the device over BLE and ANT+, has support for heart rate monitors, and even has an "exercise bike" mode which uses the trainer's buttons to change resistance... and I'm absolutely certain that there are more features which I've missed.

FortiusAnt only connects to the trainer's headunit, whereas bicipi can connect directly to the motorbrake or via the headunit.

FortiusAnt is somewhat finnicky to install, as you must install the dependencies yourself, and its dependency list is long and somewhat outdated.

FortiusAnt doesn't seem to support MyWhoosh.

That said, WouterJD is still maintaining FortiusAnt, so if you're looking for a full-featured app, go check it out.

## Useful links and credits

This project is almost entirely based on the hard work by totalreverse and WouterJD. Big thanks to them for their time and effort.

* [BLE FTMS profile](https://www.bluetooth.com/specifications/specs/fitness-machine-profile-1-0/)
* https://github.com/totalreverse/ttyT1941
* https://github.com/WouterJD/FortiusANT
* https://github.com/GoldenCheetah/GoldenCheetah/blob/80ba6ea06a6fdd5a951b028e1baf6a9810613ca0/src/Train/Fortius.h#L208
* https://github.com/WouterJD/FortiusANT/issues/171#issuecomment-748359469
* http://www.kreuzotter.de/english/espeed.htm
* https://www.gribble.org/cycling/power_v_speed.html