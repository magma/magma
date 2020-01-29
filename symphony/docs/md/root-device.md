---
id: root-device
title: Unlock Android Devices
---

## Note

This guide is tested to work on Google Pixel3 and Essential PH1, but should work for other Android phones as well. All examples and screenshots used in this guide come from Google Pixel3 and are using MacOS for the computer.

You will also need a device running _Android 4.2 or higher_ to complete this guide.


## Overview
This document will guide you through the process of side-loading `Magisk` and `TWRP` which will then allow you grant the Tech App carrier privileges.  Carrier privileges are required to do full cell scans, otherwise the app will only report cell towers that are registered to the SIM card.

`Magisk` is an open source software that can be installed on Android to root Android phones. However, with an out-of-box Android phone, it is impossible to install `Magisk` directly.

Every OEM ships phones with its own recovery system so that users can recover from system failures by, for instance, doing a factory reset. The first step in unlocking an Android phone is to unlock the bootloader, which allows us to then flash a custom recovery system called `TWRP` to replace the original recovery. `TWRP` can install packages preloaded on the phone’s sdcard. Therefore, we can preload `Magisk` on to the sdcard, flash `TWRP`, enter `TWRP` recovery and install `Magisk` from there.


## Backup (**IMPORTANT**)

Please back up any important data before proceeding further as there are risks of data loss during the process. _On some phones (including Google Pixel3), rooting the phone requires the phone to be reset_.  ALL DATA COULD BE LOST!


## Preparation (2min)

### 1. Enable developer options

If have not done so, proceed to `Settings→About` and click on `Build number` 10 times to enable developer option.

<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/mobile-app/unlock-device/developer-options.png' width=300>

### 2. Developer settings

Proceed to `Settings->System->Developer options` and enable the following

1. USB Debugging
2. OEM Unlocking (if there is no such option, then unfortunately you may not complete this guide and unlock your phone)

### 3. Connect your phone to your computer

Allow any popups that request USB debugging support from your computer.

## Software (3min)

### 1. ADB + fastboot (1min)

Download Mac platform tools through:
https://dl.google.com/android/repository/platform-tools-latest-darwin.zip

Unzip the file and place the extracted `platform-tool` folder wherever you want to. In this guide we place it under `~/Desktop` so we now have a folder at `~/Desktop/platform-tool`

### 2. TWRP (1MIN)

Visit https://twrp.me/Devices/, find the device you have and download the following files

1. `twrp-xxxx.zip`
2. `twrp-xxxx.dmg`

In this example we use Google Pixel3, so we download the following files from https://dl.twrp.me/blueline/

1. `[twrp-pixel3-installer-blueline-3.3.0-0.zip](https://dl.twrp.me/blueline/twrp-pixel3-installer-blueline-3.3.0-0.zip.html)`
2. `[twrp-3.3.0-0-blueline.img](https://dl.twrp.me/blueline/twrp-3.3.0-0-blueline.img.html)`

Note: _It is OK if your `TWRP` version is different from what we use as example._ This guide does not depend on a particular `TWRP` version.

Place both `twrp-xxxx.zip` and `twrp-xxxx.dmg` under the `platform-tool` folder you just extracted. In our example we now have the following files on our machine:

1. `~/Desktop/platform-tools/twrp-pixel3-installer-blueline-3.3.0-0.zip`
2. `~/Desktop/platform-tools/twrp-3.3.0-0-blueline.img`

### 3. magisk (1min)

Visit https://github.com/topjohnwu/Magisk/releases and find latest `Magisk` release (not `Magisk Manager`). Download `Magisk-vXXX.zip` and place it under the `platform-tool` folder you just extracted.

In our example we now have the following file on our machine:
`~/Desktop/platform-tools/Magisk-v19.1.zip`

Note: _It is OK if you download a newer version of `Magisk` than what we use as example._ This guide does not depend on a particular `Magisk` version.


## Unlock the Phone (10-20min)

_Any command used in this section starts with `$`. Do not include the `$` when you actually type the command. It is just used to indicate the start of a command._
For instance, if the guide asks you to type `$ cd ~/Desktop`, you only need to type `cd ~/Desktop` on your computer.

### 1. Setup terminal (30sec)

1. Open terminal on your Mac. One way to do so is to open `Spotlight` and search for `Terminal`
2. Once inside the terminal, go to the `platform-tool` folder by typing the following command
    `$ cd <path-to>/platform-tools
    `In our example we type `$ cd ~/Desktop/platform-tools`
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/mobile-app/unlock-device/find-terminal.png' width=600>

### 2. Enter `fastboot` mode (30sec)

Run the following command in terminal.

```$ adb reboot bootloader```

Wait for the phone to restart and enter fastboot mode.

### 3. Unlock the bootloader (1min-10min)

Run the following command in terminal.

```$ fastboot flashing unlock```

Then use volume up/down and power button to confirm `Unlock bootloader` if asked to.

If the phone is reset after unlocking the bootloader, repeat all steps in the Preparation section. You do not need to do anything if the phone is only rebooted but not reset.

### 4. Push custom recoveries and magisk (30sec)

Go back to your terminal window and run the following commands:
```
$ adb push ./twrp-xxxx.zip /sdcard
$ adb push ./Magisk-vXXX.zip /sdcard
```
where `twrp-xxxx.zip` is downloaded in Software section step2 and `Magisk-vXXX.zip` is downloaded in Software section step3.

In our example we run:
```
$ adb push ./twrp-pixel3-installer-blueline-3.3.0-0.zip /sdcard
$ adb push ./Magisk-v19.1.zip /sdcard
```
### 5. Set BOOT slot (30sec)

First run the following command and wait for the phone to enter fastboot mode again:

```$ adb reboot bootloader```

Once in fastboot mode, run the following command:

```$ fastboot getvar current-slot```

An example output looks like the following:

```
current-slot: a
Finished. Total time: 0.080s
```

If your current slot is a, run the following commands:

```
$ fastboot flash boot_b ./twrp-xxxx.img
$ fastboot --set-active=b
```
If your current slot is b, run the following commands:
```
$ fastboot flash boot_a ./twrp-xxxx.img
$ fastboot --set-active=a
```
Here `twrp-xxxx.img` is downloaded in Software section step 2. In our example the file is
`./twrp-3.3.0-0-blueline.img`

### 6. Enter recovery mode (30sec)

Use the volume up/down key to select `Recovery Mode` and press the power button to enter recovery mode. Expect to see the following screen, which means you have entered TWRP recovery.

<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/mobile-app/unlock-device/TWRP-recovery-mode.png' width=300>

### 7. Install TWRP (30sec)

Follow the steps below to install TWRP.

1. Press `Keep Read Only` on TWRP welcome screen.
2. Press `Install` at the upper left corner.
3. Find the file `twrp-xxxx.zip` and press on it.
4. Swipe to confirm flash.

### 8. Reboot the phone (1min)

1. Press home button to go back to home screen of TWRP.
2. Press `Reboot`.
3. If your current-slot in step5 was `a`, press on `Slot A`. Otherwise press on `Slot B`.
4. Press home button and then press `Reboot` again.
5. Press on `System` at upper left corner and then choose `Do Not Install` to reboot the phone.

### 9. Enter bootloader (1min)

Once the phone is rebooted, type the following command to reenter bootloader mode.

```
$ adb reboot bootloader
```
After the phone reboots, use volume up/down button to select `Recovery Mode` and press power button to enter.

### 10. Install magisk (1min)

You should now see the TWRP home screen again. Follow steps below to install Magisk.

1. Press `Install` at upper left corner.
2. Find file `Magisk-vXXX.zip` and press on it. In our example, the file is `Magisk-v19.1.zip`
3. Swipe to confirm flash.

### 11. Reboot the phone (30sec)

1. Press home button to go back to home screen of TWRP.
2. Press `Reboot` at the lower right corner.
3. Press `System` and then `Do Not Install` to reboot the phone.

### 12. Confirm (30sec)

Wait for 30 seconds after the phone reboots. Check if `Magisk` is installed. If so, then your phone is successfully rooted.

<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/mobile-app/unlock-device/confirm-magisk-install.png' width=300>

## Additional Setup: Grant System Access to the Mobile App (5min)

_This setup only needs to be completed once._

### 1. Install root file explorer (1min)

Use Google Play to install a root file explorer app. The example we use is `File Explorer Root Browser`

<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/mobile-app/unlock-device/root-file-explorer.png' width=300>

### 2. Move permissions file to laptop (1min)

1. Use file explorer to navigate to `/etc/permissions` and locate `privapp-permissions-platform.xml`. (_The file may be named slightly different on different phones_).
2. Use `Share` inside file explore to move the file to your laptop.

### 3. Register the mobile app as privileged app (1min)

Use any text editor to add the following lines right above the last line in `privapp-permissions-platform.xml`:
```
<privapp-permissions package="cloud.thesymphony">
    <permission name="android.permission.MODIFY_PHONE_STATE"/>
</privapp-permissions>
```
The end of file should now look like this:

<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/mobile-app/unlock-device/permission-file-snippet.png' width=600>

### 4. Copy permission file back (2min)

1. Run the following command to push modified permission file back to the phone.
```
$ adb push <path-to>/privapp-permissions-platform.xml /sdcard
```
2. Use root file explorer to move modified file back to `/etc/permissions/`privapp-permissions-platform.xml``
3. Select file `/etc/permissions/privapp-permissions-platform.xml`, go to `properties` in the file explorer. Modify the permission to `-rw-r--r--` and press `Apply`.

<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/mobile-app/unlock-device/file-permissions.png' width=300>

### 5. Reboot the phone (30sec)

If your phone fails to start and enters TWRP again, you may have modified the file incorrectly or set the permission incorrectly in step4. To recover:

1. Press `Advanced` on TWRP home screen and select `fix boot loop`
2. After the phone boots, visit https://developers.google.com/android/ota and follow the instructions to do a factory reset.

## Additional Section: Grant System Access to the Technician App (1min)

This section needs to be performed _every time a new version of the app is installed_.

### 1. Transfer app to system space

1. Use root file explorer to navigate to `/data/app` and locate a folder starting with `cloud.thesymphony`.
2. Cut and paste the folder to `/system/priv-app`

### 2. Reboot the phone

If your phone fails to start and enters TWRP again, you may have modified the permission file incorrectly in step 4 of previous section. Please follow step 5 of the previous section to fix.
