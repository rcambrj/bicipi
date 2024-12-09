# bicipi

status: incomplete

**bicipi**: pronounced _bee-see-pie_, is the fusion of the Spanish word for bike + Raspberry Pi

Connects to a USB Tacx T1941 trainer and exposes it via Bluetooth Low Energy for use with Zwift, TrainerRoad or MyWhoosh.

## Objectives

* Lightweight and fast:
	* present the device on BLE just like a modern trainer would (no GUI)
* Easy to install:
	* have a downloadable + flashable image for a range of Raspberry Pis
* Support trainers I have access to:
	* Tacx Fortius T1941
* Support popular apps:
	* Zwift, TrainerRoad, MyWhoosh

## FortiusAnt

FortiusAnt by WouterJD is a comprehensive software package. It connects to a wide range of hardware, provides a GUI, exposes the device over BLE and ANT+, has support for heart rate monitors, and even has an "exercise bike" mode which uses the trainer's buttons to change resistance... and I'm absolutely certain that there are more features which I've missed.

FortiusAnt connects to the trainer's headset via USB, whereas bicipi connects directly to the motorbrake, bypassing the headset. Connecting via USB is trivial and bicipi requires some tools and knowledge in order to make the connection.

FortiusAnt's installation is somewhat involved and finnicky, as its dependency list is long and somewhat outdated.

That said, WouterJD is still maintaining FortiusAnt, so if you're looking for a full-featured app, go check it out.

## Useful links and credits

* https://github.com/totalreverse/ttyT1941
* https://github.com/WouterJD/FortiusANT
* https://www.bluetooth.com/specifications/specs/fitness-machine-profile-1-0/