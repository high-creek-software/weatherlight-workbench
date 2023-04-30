# Weatherlight Workbench

## Packaging
Fyne apps can be packaged for many different platforms.

### Android
- Install android sdk/ndk
- Set ANDROID_NDK_HOME
- cd cmd/weatherlightworkbench
- fyne package -os android -appID io.highcreeksoftware.weatherlightworkbench --icon ../../internal/platform/icons/app_icon.png
- adb install weatherlightworkbench.apk