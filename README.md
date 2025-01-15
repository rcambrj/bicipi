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

Download the zip for your pi from the releases page and use [etcher](https://etcher.balena.io/) to burn it to a USB stick. Insert the USB stick into your pi and give it power. After a minute or two, connect the Tacx via USB. After a few seconds you should be able to kick the pedal to begin the calibration process.

The calibration process by default warms the motor and tyre up for 5 minutes. If the motor is still spinning after 10 minutes, something is wrong.

To find out what's going on, you may connect your PC directly to the Raspberry Pi via a network cable, configure your PC with address `192.168.24.123` and subnet `255.255.255.0`, and run `ssh bicipi@192.168.24.24`. If you're prompted for a password, use `password`. You can now scan the logs for the `bicipi` systemd service with `journalctl -xfeu bicipi` and restart it with `sudo systemctl restart bicipi`.

> [!IMPORTANT]
> DO NOT connect this Raspberry Pi to your regular network. The password is easily guessed.

Known issues:

* it must be an USB stick - SD cards won't boot right now for some reason
* if the Tacx is connected via USB at boot time, it will hang until it's disconnected
* there is no way to configure any settings right now


## CLI installation and usage

Dev tooling and dependencies are declared in the nix devshell, for which `nix` is required and `direnv` is recommended.

The tool is a simple Go executable, you can run it with then: `go run .`. It only runs on Linux, it spits out some errors on MacOS.

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

Raspberry Pis have limited UARTs, so the direct-to-motorbrake serial connection is discouraged on these boards. Feel free to connect a ttyUSB and use that though.

See @totalreverse's [diagram](https://github.com/totalreverse/ttyT1941/wiki#inside-the-t1902-and-t1942) of the `Brake backside socket` for how to connect.

Remember to put [a pull-up resistor](https://github.com/totalreverse/ttyT1941/issues/7#issuecomment-619587334) on the RX line (motorbrake TX).

## FortiusAnt

FortiusAnt by WouterJD is a comprehensive software package. It connects to a wide range of hardware, provides a GUI, exposes the device over BLE and ANT+, has support for heart rate monitors, and even has an "exercise bike" mode which uses the trainer's buttons to change resistance... and I'm absolutely certain that there are more features which I've missed.

FortiusAnt connects to the trainer's headset via USB, whereas bicipi connects directly to the motorbrake, bypassing the headset. Connecting via USB is trivial and bicipi requires some tools and knowledge in order to make the connection.

FortiusAnt's installation is somewhat involved and finnicky, as its dependency list is long and somewhat outdated.

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