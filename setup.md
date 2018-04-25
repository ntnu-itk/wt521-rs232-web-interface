# WT521 setup

1. connect to WT521 with RS-232
  - e.g. screen /dev/ttyS0 9600,n,cs8
2. power on WT521 and, within the first 5 seconds, send "open" to start config.
3. see the manual's page 47 (PDF's page 49) for commands
4. send "reset" to reset WT521 and load your settings

## Example setup

```
# Get current settings
getset
# Set measurement #1 to be of the MWV message format (see manual's page 49)
setmes 1 type MWV
# Make the MWV messages be sent every 172.8 seconds => 500 samples/day
setmes 1 interval 172.8
# Adjust wind direction offset it your unit is not physically oriented with the
# WAV151 direction sensor towards North ([-180, 180], degrees)
setdir offset 0
# Configure RS-232 interface to use desired baud rate (1200 is default)
setcom 0 baudrate 1200
```
## My setup

This is the output of the "getset" command on the WT521 that I tested this
software with:

```
>getset
SETDEV  ID=A
SETSPD  AVERAGE=3.00
SETDIR  AVERAGE=3.00 
SETMEA
 WA151  ACTIVE=1
 WMS301 ACTIVE=0
 WMS302 ACTIVE=0
 HMP45  ACTIVE=0
 PT100  ACTIVE=0
SETMES
 1      TYPE=MWV COM=0 INTERVAL=172.80
 2      TYPE=NONE
 3      TYPE=NONE
 4      TYPE=NONE
SETCOM
 0      BAUDRATE=1200 DATABITS=8 STOPBITS=1 PARITY=NONE WIRES=2
 1      DATABITS=7 STOPBITS=2 PARITY=NONE
        CCITT=V.21 MULTIDROP=0 ORIGINATE=1
SETHEA  ACTIVE=1
```