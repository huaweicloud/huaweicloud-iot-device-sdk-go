English | [简体中文](./README_CN.md)

# huaweicloud-iot-device-sdk-go

# 0. Version update instructions
| Version | Change Type | Description |
|:-------|:-|:-------------------------------------------------|
| v1.0.0 | New features | Provides the ability to connect to the Huawei Cloud IoT platform to facilitate users to implement business scenarios such as secure access, device management, data collection, command issuance, device provisioning, and client-side rules |   

# 1. Preface
huaweicloud-iot-device-sdk-go provides a Go version of the SDK for devices to access the Huawei Cloud IoT platform. It provides communication capabilities between devices and platforms, as well as advanced services such as device services, gateway services, and OTA. It also provides various The scene provides rich demo codes. IoT device developers can use the SDK to greatly simplify development complexity and quickly access the platform.

This article uses examples to describe how the SDK helps devices quickly connect to the Huawei IoT platform using the MQTT protocol.

Huawei Cloud official website: https://www.huaweicloud.com/

Click "Console" in the upper right corner of the Huawei Cloud official website to enter the management console, and search for "IoTDA" at the top of the page to enter the device access service console.

# 2.SDK Introduction
## 2.1 Introduction to SDK functions
The SDK is oriented to embedded terminal devices with strong computing and storage capabilities. Developers can achieve uplink and downlink communication between the device and the IoT platform by calling the SDK interface. The functions currently supported by the SDK are:
* Supports device messages, attribute reporting, attribute reading and writing, and command issuance
*Support OTA upgrade
* Supports two device authentication methods: password authentication and certificate authentication
* Support device shadow query
* Support gateway and sub-device management
* Support custom topics
* Support file upload/download
*Support end-side rules
* Support equipment distribution

## 2.2SDK directory structure
<table>
  <tr>
    <td>Directory structure</td>
    <td>The Table of Contents</td>
    <td>Description</td>
  </tr>
 <tr>
   <td rowspan=9>iot</td>
   <td> callback</td>
   <td> Client Callback Function</td>
  </tr>
  <tr>
   <td>client</td>
   <td>Device Client</td>
  </tr>
  <tr>
   <td>config</td>
   <td>Client Configuration</td>
  </tr>
  <tr>
   <td>constants</td>
   <td>constant package</td>
  </tr>
  <tr>
   <td>device</td>
   <td>Directly Connected Device Client</td>
  </tr>
  <tr>
   <td>file</td>
   <td>File upload and download</td>
  </tr>
  <tr>
   <td>gateway</td>
   <td>Gateway Device Client</td>
  </tr>
  <tr>
   <td>model</td>
   <td>Structure package</td>
  </tr>
  <tr>
   <td>rule</td>
   <td>Device-side rule package</td>
  </tr>
<tr>
   <td rowspan=15>sample</td>
   <td> bs</td>
   <td> Device bootstrap demo</td>
  </tr>
  <tr>
   <td>command</td>
   <td>Command demo</td>
  </tr>
  <tr>
   <td>file</td>
   <td>File upload and download demo</td>
  </tr>
  <tr>
   <td>device</td>
   <td>Directly Connected Device Demo</td>
  </tr>
  <tr>
   <td>file</td>
   <td>File upload and download Demo</td>
  </tr>
  <tr>
   <td>gateway</td>
   <td>Gateway management demo</td>
  </tr>
  <tr>
   <td>log</td>
   <td>Device Log demo</td>
  </tr>
  <tr>
   <td>message</td>
   <td>Message reporting and delivery demo</td>
  </tr>
  <tr>
   <td>ota</td>
   <td>ota upgrade demo</td>
  </tr>
  <tr>
   <td>properties</td>
   <td>device properties demo</td>
  </tr>
  <tr>
   <td>rule</td>
   <td>Device rule demo</td>
  </tr>
  <tr>
   <td>test_model</td>
   <td>Test data structure</td>
  </tr>
  <tr>
   <td>test_util</td>
   <td>Test tool</td>
  </tr>
  <tr>
   <td>test_sync</td>
   <td>Time synchronization demo</td>
  </tr>
</table>

# 3. Preparation work
* Go 3.18 installed
* Relevant dependencies have been installed according to go.mod

# 4.SDK function
## 4.1 Upload product model and register device
In order to facilitate the experience, we provide a smoke sensor product model. The smoke sensor will report smoke value, temperature, humidity, smoke alarm, and also supports ring alarm commands.
Take the smoke sensor as an example to experience functions such as message reporting, attribute reporting, and command response.

* Visit [Device Access Service](https://www.huaweicloud.com/product/iothub.html) and click "Use Now" to enter the device access console.

* Click "Access Information" to view the MQTT device access address and save the address.

   ![](.\doc\figure_en\get_access_address_en.png)	


* Select "Product" in the device access console, click "Create Product" in the upper right corner, and in the pop-up page, fill in the "Product Name", "Protocol Type", "Data Format", "Manufacturer Name", " Industry", "Equipment Type" and other information, and then click "Create Now" in the lower right corner.

   - Select "MQTT" as the protocol type;

   - Select "JSON" as the data format.
 
   ![](.\doc\figure_en\upload_profile_2_en.png)

* After the product is successfully created, click "Details" to enter the product details. On the function definition page, click "Upload Model File" to upload the smoke detector product model [smokeDetector](https://iot-developer.obs.cn-north -4.myhuaweicloud.com/smokeDetector.zip).
    The generated product model is shown in the figure below.

    ![](.\doc\figure_en\upload_profile_2_1_en.png)

* In the left navigation bar, select "Devices > All Devices", click "Register Device" in the upper right corner, in the pop-up page, fill in the registration device parameters, and then click "OK".

   ![](.\doc\figure_en\upload_profile_3_en.png)

* After the device is successfully registered, the device identification code, device ID, and key are saved.

## 4.2 Online debugging tools
In the left navigation bar of the console, select "Monitoring and Operation > Online Debugging" to enter the online debugging page.
The page has command issuing and message tracking functions.

* Click "Select Device" in the upper right corner of the page to select the registered device

* Clicking "IoT Platform" will display message tracking

* Click "Send" in the lower right corner of the page to send commands to the device


## 4.3 Device initialization

* Create device
* Create device.

   When a device is connected to the platform, the IoT platform provides two authentication methods: key and certificate.

   * If you use port 1883 to access the platform through key authentication, you need to write the obtained device ID and key.

   ```go
    //Create a device and connect to the platform
	authConfig := &amp;config2.ConnectAuthConfig{
		Id:       "{your device id}",
		Servers:  "mqtt://{access_address}:1883",
		Password: "your device secret",
	}
	mqttDevice := device2.NewMqttDevice(authConfig)
   ```

   * If you use port 8883 to access the platform through key authentication (recommended, all SDK demos access the platform through this method), you need to write the obtained device ID, key and preset CA certificate.
   Preset certificate:/samples/resources/root.pem
  
   ```go
    authConfig := &amp;config.ConnectAuthConfig{
		Id:       "{your device id}",
		Servers:  "mqtts://{access_address}:8883",
		Password: "your device secret",
		ServerCaPath: "./resources/root.pem",
	}
	mqttDevice := device.NewMqttDevice(authConfig)
   ```
   * If you use port 8883 and access the platform through X509 certificate authentication, you need to write the obtained device ID, certificate information and pre-made CA certificate. For more X509 certificate access, please refer to [X509 Certificate Access](https://support.huaweicloud.com/bestpractice-iothub/iot_bp_0077.html)
     Preset certificate:/samples/resources/root.pem
   
    ```go
     authConfig := &amp;config.ConnectAuthConfig{
		 Id:       "{your device id}",
		 Servers:  "mqtts://{access_address}:8883",
		 AuthType:        constants.AuthTypeX509,
		 ServerCaPath: "./resources/root.pem",
         CertFilePath: "your device cert path",
		 CertKeyFilePath: "your device cert key path",
	 }
	 mqttDevice := device.NewMqttDevice(authConfig)
    ```

* Call the connect interface to establish a connection. This interface is a blocking call and will return true if the connection is successfully established.

   ```go
        initResult := mqttDevice.Connect()
   ```

* After the connection is successful, communication between the device and the platform begins. Obtain the device client through the Client attribute of MqttDevice. The client provides communication interfaces such as messages, properties, and commands.
For example:
   ```go
        mqttDevice.Client.AddCommandHandler(func(command model.Command) (bool, interface{}) {
		    fmt.Println("I get command from platform")
		    return true, map[string]interface{}{
			    "cost_time": 12,
		    }
	    })
        mqttDevice.SendMessage(message)
   ```

* For detailed information about the MqttDevice class, see iot/device/mqtt_device.go

If the connection is successful, the "Message Tracking" on the online debugging page will display:

![](.\doc\figure_en\init_1_en.png)

The running log is:

![](.\doc\figure_en\init_2_en.png)



## 4.4 Command issuance

/samples/command/platform_command.go is an example of processing platform command issuance.
Set up a command listener to receive commands issued by the platform. In the callback interface, the commands will be processed and the responses will be reported.

The following code implements the CommandHandler function as a command listener, namely:

```go
    mqttDevice.Client.AddCommandHandler(func(command model.Command) (bool, interface{}) {
		glog.Infof("I get command from platform")
		glog.Infof("command device id is %s", command.ObjectDeviceId)
		glog.Infof("command name is %s", command.CommandName)
		glog.Infof("command serviceId is %s", command.ServiceId)
		glog.Infof("command params is %v", command.Paras)
		return true, map[string]interface{}{
			"cost_time": 12,
		}
	})
```

The CommandHandler function will be automatically called when the device receives the command.
The example prints the content of the command in the CommandHandler function and returns the response to the platform.

Execute the main method and issue commands to the device on the "Online Debugging" page. The code will produce the following output:

![](.\doc\figure_en\command_1_en.png)

At the same time, the device's response to the command can be found in the "Message Trace" of "Online Debugging".

![](.\doc\figure_en\command_2_en.png)


## 4.5 Platform message delivery/device message reporting
Message delivery refers to the platform delivering messages to the device. Message reporting refers to the device reporting messages to the platform. For more device message information, please refer to [Device Message Document](https://support.huaweicloud.com/usermanual-iothub/iot_01_0322.html)

### 4.5.1 Message reporting
/samples/message/message_demo.go is an example of message reporting.
```go
    """ create device code here """
    //Send message to platform
    message := model.Message{
		Topic:   "{custom topic}",
		Payload: "first message",
	}

	for {
		sendMsgResult := device.SendMessage(message)
		glog.Infof("send message %v", sendMsgResult)
		time.Sleep(1 * time.Second)
	}
```
In the above code, you can report the message to a custom topic through the SendMessage method, or you can not fill in the topic in the message, and the message will be reported to the default topic of the platform. If sent successfully, you can see on the "Online Debugging" page:

![](.\doc\figure_en\message_1_en.png)

### 4.5.2 Message delivery
/samples/message/message_demo.go is an example of message delivery.
```go
    < create device code here ... >

    // Register the callback for messages sent by the platform. When receiving the message sent by the platform, call this callback.
	// Supports registering multiple callbacks and calling them in the order of registration
	device.Client.AddMessageHandler(func(message string) bool {
		glog.Infof("first callback called : %s", message)
		return true
	})
    // You can implement some custom logic here after establishing a link with the platform, such as subscribing to some custom topics after establishing a link.
	device.Client.ConnectHandler = func(client mqtt.Client) {
		// Use a custom topic to receive messages sent by the platform
		device.Client.SubscribeCustomizeTopic("testdevicetopic", func(message string) bool {
			glog.Infof("first callback called %s", message)
			return true
		})
	}
    device.Connect()

```

In the above code, the ConnectHandler can subscribe to messages delivered by a custom topic through SubscribeCustomizeTopic after the platform establishes a link. If messages delivered by a custom topic are not used, messages delivered by the platform's default topic are accepted through the AddMessageHandler method. After executing the main function, you can use the platform to deliver messages. The code will produce the following output:
![](.\doc\figure_en\message_2_en.png)

## 4.6 Properties reporting/setting
Properties reporting refers to the device reporting the current attribute values ​​to the platform. Property settings refer to the platform setting property values ​​of the device.
/samples/properties/device_properties.go is an example of property reporting/setting.
### 4.6.1 Properties reporting
Used by the device to report properties data to the platform in the format defined in the product model. The platform will assign the reported data to the device shadow data.


   ```go
    < create device code here ... >
       //Device reporting properties
	props := model.DevicePropertyEntry{
		ServiceId: "smokeDetector",
		EventTime: iot.GetEventTimeStamp(),
		Properties: test_model.DemoProperties{
			Temperature: 28,
		},
	}

	var content []model.DevicePropertyEntry
	content = append(content, props)
	services := model.DeviceProperties{
		Services: content,
	}
	device.ReportProperties(services)
   ```

The above code will report the temperature attribute.
If the report is successful, the "Online Debugging" page will display:

![](.\doc\figure_en\properties_1_en.png)

Select "Devices > All Devices" in the left navigation bar, select the registered device to view, and you can see the attribute values ​​just reported in "Device Shadow".

![](.\doc\figure_en\properties_2_en.png)

### 4.6.2 Platform settings device properties
If the device is set as a property listener through the AddPropertiesSetHandler and SetPropertyQueryHandler methods, that is:

```go
//Register platform setting attribute callback. When the application sets device attributes through the API, this callback will be called. It supports registering multiple callbacks.
	device.Client.AddPropertiesSetHandler(func(propertiesSetRequest model.DevicePropertyDownRequest) bool {
		fmt.Println("I get property set command")
		fmt.Printf("request is %s", iot.Interface2JsonString(propertiesSetRequest))
		return true
	})
    // Register the platform query device attribute callback. This callback is called when the platform queries device attributes. Only one callback is supported.
	device.Client.SetPropertyQueryHandler(func(query model.DevicePropertyQueryRequest) model.DevicePropertyEntry {
		return model.DevicePropertyEntry{
			ServiceId: "smokeDetector",
			Properties: test_model.DemoProperties{
				Temperature: 27,
			},
			EventTime: "2024-05-28 14:23:24",
		}
	})
```

Then when the device receives a property read and write request, it will automatically call the PropertiesSetHandler or PropertyQueryHandler method in the listener.
Among them, the PropertiesSetHandler method handles writing properties, and the PropertyQueryHandler method handles reading attributes.
In most scenarios, users can read the device shadow directly from the platform, so the PropertyQueryHandler method does not need to be implemented.
But if you need to support real-time reading of properties from the device, you need to implement this method.
The example prints the content of the property settings in the PropertiesSetHandler method and returns the response to the platform.

```go

def run():
    < create device code here ... >
    
    //Register platform setting attribute callback. When the application sets device attributes through the API, this callback will be called. It supports registering multiple callbacks.
	device.Client.AddPropertiesSetHandler(func(propertiesSetRequest model.DevicePropertyDownRequest) bool {
		glog.Infof("I get property set command")
		glog.Infof("request is %s", iot.Interface2JsonString(propertiesSetRequest))
		return true
	})
    device.Connect()
	glog.Infof("device connected: %v\n", device.IsConnected())

```

At "Device Shadow", click "Property Configuration" to set the expected value of the property.
If the set expected value is different from the value reported by the device, the platform will automatically send the expected value to the device when the device goes online. (That is, the platform sets device properties)

![](.\doc\figure_en\properties_3_en.png)

Run the above run function and get:

![](.\doc\figure_en\properties_4_en.png)

## 4.7 Device Shadow
Used by the device to obtain device shadow data from the platform. The device can obtain the platform device shadow data to synchronize the device attribute values, thereby completing the modification of the device attribute values.

/samples/properties/device_properties.go also contains device acquisition platform device shadow data.

* The device requests to obtain the device shadow data of the platform.

   ```go
    # Receive platform downlink response
    device.Client.AddDeviceShadowQueryResponseHandler(func(response model.DeviceShadowQueryResponse) {
		shadow := response.Shadow
		glog.Infof("receive shadow msg from device.")
		glog.Infof("on_shadow_get device_id:  %s", response.ObjectDeviceId)
		glog.Infof("shadow service_id: %s", shadow[0].ServiceId)
		glog.Infof("shadow desired properties: %v", shadow[0].Desired)
		glog.Infof("shadow reported: %v", shadow[0].Reported)
	})

    device.Connect()
	glog.Infof("device connected: %v\n", device.IsConnected())

    //Device query device shadow data
	device.QueryDeviceShadow(model.DevicePropertyQueryRequest{
		ServiceId: "smokeDetector",
	})
   ```
## 4.8 OTA upgrade
An example of OTA upgrade is implemented in /samples/ota/ota_demo.go, as shown in the following code.

```go
   < create device code here ... >

    // OTA listener settings, report SDK version number
    device.Client.SwFwVersionReporter = func() (string, string) {
		return "v1.0", "v1.0"
	}

	// upgradeType 0: software upgrade 1: firmware upgrade
	device.Client.DeviceUpgradeHandler = func(upgradeType byte, info model.UpgradeInfo) model.UpgradeProgress {
		glog.Infof("begin to handle upgrade process")
		upgradeProcess := model.UpgradeProgress{}
		currentPath, err := os.Getwd()
		if err != nil {
			glog.Warningf("get executable path failed. err: %s", err.Error())
			upgradeProcess.ResultCode = 255
			upgradeProcess.Description = "get executable path failed."
			return upgradeProcess
		}
        // Save path of software and firmware download package
		currentPaths = currentPath + "\\download\\ota.txt"
		downloadFlag := file.CreateHttpClient().DownloadFile(currentPath, info.Url, info.AccessToken)
		if !downloadFlag {
			glog.Errorf("down load file { %s } failed", currentPath)
			upgradeProcess.ResultCode = 10
			upgradeProcess.Description = "down load ota package failed."
			return upgradeProcess
		}
		glog.Infof("download file success.")
		//checkPackage() Verifies the downloaded upgrade package
		//installPackage() installs the upgrade package
		upgradeProcess.ResultCode = 0
		upgradeProcess.Version = info.Version
		upgradeProcess.Progress = 100
		return upgradeProcess
	}

    connect := device.Connect()
	glog.Infof("connect result : %v", connect)
```
What the user needs to implement is the DeviceUpgradeHandler listener. This example is an example of listener implementation.

### 4.8.1 How to perform OTA upgrade

1. Firmware upgrade. Refer to [Firmware Upgrade](https://support.huaweicloud.com/usermanual-iothub/iot_01_0027.html)

2. Software upgrade. Refer to [Software Upgrade](https://support.huaweicloud.com/usermanual-iothub/iot_01_0047.html)

## 4.9 File upload/download
An example of file upload/download is implemented in /samples/file/upload_file.go.

```go
    < create device code here ... >
    
    fileName := "test_upload.txt"
	uploadFilePath := currentPath + "\\download\\test_upload.txt"
    //Upload file
	device.UploadFile(fileName, uploadFilePath)

	time.Sleep(10 * time.Second)
    // Download file
	downloadFilePath := currentPath + "\\download\\test_download.txt"
	device.DownloadFile(fileName, downloadFilePath)
```

File upload/download process reference [File Upload](https://support.huaweicloud.com/usermanual-iothub/iot_01_0033.html)

* Configure OBS storage in the console.
   
   ![](.\doc\figure_en\obs_config_en.png)

* Preset upload files. The file to be uploaded in the above example is /iot_device_demo/filemanage/download/upload_test.txt.
   In the file download part, download the uploaded upload_test.txt and save it to /iot_device_demo/filemanage/download/download.txt.

* Execute the above example to see the storage results on OBS.
   
   ![](.\doc\figure_en\obs_object_en.png)

## 4.10 Device time synchronization
An example of device time synchronization is implemented in /samples/time_sync/time_sync_demo.go.

```go
   < create device code here ... >

    // Set up time synchronization service
    device.Client.SyncTimeResponseHandler = func(deviceSendTime, serverRecvTime, serverSendTime int64) {
		deviceRecvTime := time.Now().UnixNano() / int64(time.Millisecond)
		now := (serverRecvTime + serverSendTime + deviceRecvTime + deviceSendTime) / 2
		glog.Infof("now is %d", now)
	}

    connect := device.Connect()
	glog.Infof("connect result : %v", connect)

    # Request time synchronization
    sync := device.RequestTimeSync()
    glog.Infof("sync time result: %v", sync)
	time.Sleep(10 * time.Second)
```
Users can implement the SyncTimeResponseHandler method themselves.
The SyncTimeResponseHandler method in /samples/time_sync/time_sync_demo.go is an example of a listener implementation. Assuming that the device-side time received by the device is device_recv_time, the device calculates its own accurate time as:
(server_recv_time + server_send_time + device_recv_time - device_send_time) / 2

## 4.11 Gateway and sub-device management
For this function, please refer to [Gateway and Sub-Device](https://support.huaweicloud.com/usermanual-iothub/iot_01_0052.html)

The demo code for gateway and sub-device management is under /samples/gateway/gateway_demo.go. This demo demonstrates how to use a gateway device to communicate with the platform.

This demo can demonstrate:
1. Gateway synchronization sub-device list. When the gateway device is not online, the platform cannot notify the gateway device of new and deleted sub-device information in a timely manner.
2. When the gateway device goes offline and comes online again, the platform will notify the newly added/deleted sub-device information.
3. The gateway updates the status of the sub-device. The gateway notifies the platform that the status of the sub-device is "ONLINE".
4. The sub-device reports the message to the platform through the gateway.
5. The platform issues commands to the sub-devices.
6. Gateway add/delete sub-device request

### 4.11.1 Create gateway device
* Consistent with direct connection devices, port 8883 can be used to access the platform
```go
    authConfig := &amp;config.ConnectAuthConfig{
		Id:           "your device id",
		Servers:      "mqtts://{your access ip}:8883",
		Password:     "your device key",
		ServerCaPath: "./resources/root.pem",
	}
	mqttDevice := gateway.NewMqttGatewayDevice(authConfig)
	connect := mqttDevice.Connect()
	if !connect {
		return nil
	}
```
### 4.11.2 The platform notifies the addition and deletion of gateway sub-devices
In samples/gateway/gateway_demo.go, the platformNotifySubDeviceAdd() and platformNotifySubDeviceDelete() methods demonstrate the platform notification gateway sub-device addition and deletion functions.
```go
    < create device code here ... >
    //The gateway receives the callback for adding a sub-device
	gatewayDevice.Client.SubDevicesAddHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device add version: %s", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("add sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
    / The gateway receives a callback to delete the sub-device
	gatewayDevice.Client.SubDevicesDeleteHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device delete version: %s", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("delete sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
```
Users can implement the SubDevicesAddHandler and SubDevicesDeleteHandler methods by themselves. This example provides a default implementation. When the platform creates and deletes a sub-device, the gateway will receive a notification and print out the sub-device information.

### 4.11.3 Gateway synchronization sub-device list
In samples/gateway/gateway_demo.go, the syncSubDevices() method demonstrates the gateway synchronization sub-device list function.
```go
    //Synchronize the response of new sub-device
	gatewayDevice.Client.SubDevicesAddHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device add version: %s", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("add sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
	// Synchronously delete sub-device response
	gatewayDevice.Client.SubDevicesDeleteHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device delete version: %s", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("delete sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
    gatewayDevice.SyncAllVersionSubDevices()
```
Users can implement the SubDevicesAddHandler and SubDevicesDeleteHandler methods by themselves. This example provides a default implementation. When the device sends a synchronization sub-device request, the platform will send the sub-device information that needs to be deleted and added to the device.

### 4.11.4 Gateway updates sub-device status
In samples/gateway/gateway_demo.go, the updateSubDeviceStats() method demonstrates the gateway’s function of updating sub-device status.
```go
//Update sub-device status request response
	gatewayDevice.Client.SubDeviceStatusRespHandler = func(response model.SubDeviceStatusResp) {
		if len(response.SuccessfulDevices) > 0 {
			glog.Infof("success update device status.")
			for _, sucessDevice := range response.SuccessfulDevices {
				glog.Infof("update device : %s status: %s", sucessDevice.DeviceId, sucessDevice.Status)
			}
		}
		if len(response.FailedDevices) > 0 {
			glog.Infof("failed to update device status")
			for _, failedDevice := range response.FailedDevices {
				glog.Infof("failed to update device status. deviceId : %s, errCode: %s, errMsg: %s",
					failedDevice.DeviceId, failedDevice.ErrorCode, failedDevice.ErrorMsg)
			}
		}
	}
    status := model.DeviceStatus{
		DeviceId: "xxxx",
		Status:   "ONLINE",
	}
	statusInfos := make([]model.DeviceStatus, 1)
	statusInfos = append(statusInfos, status)
	subDeviceStatus := model.SubDevicesStatus{
		DeviceStatuses: statusInfos,
	}
    //Update sub-device status
    gatewayDevice.UpdateSubDeviceState(subDeviceStatus)
```
Users can implement the SubDeviceStatusRespHandler method by themselves. This example provides a default implementation. After the gateway sends a request to update the sub-device status to the platform, the platform will notify the gateway of successful devices and failed devices after receiving the request to update the device status.

### 4.11.5 Add sub-device to gateway
In samples/gateway/gateway_demo.go, the gatewayAddSubDevice() method demonstrates the gateway’s function of adding sub-devices.
```go
    //Add sub-device request response
	gatewayDevice.Client.SubDevicesAddHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sub device add version: %s", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("add sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
	deviceInfo := model.DeviceInfo{
		DeviceId:       "xxxxx",
		NodeId:         "xxxxxx",
		Name:           "xxxx",
		ParentDeviceId: "xxxxx",
		Description:    "xxxxx",
		ProductId:      "xxxxx",
	}
	infos := make([]model.DeviceInfo, 1)
	infos = append(infos, deviceInfo)
	gatewayDevice.AddSubDevices(infos)
```
Users can implement the SubDevicesAddHandler method by themselves. This example provides a default implementation. After the gateway sends a sub-device addition request to the platform, the platform will notify the gateway of the new sub-device after adding the sub-device.

### 4.11.6 Gateway deletes sub-device
In samples/gateway/gateway_demo.go, the gatewayDeleteDevice() method demonstrates the gateway's function of deleting sub-devices.
```go
    // Delete child device request response
	gatewayDevice.Client.SubDevicesDeleteHandler = func(devices model.SubDeviceInfo) {
		glog.Infof("sync sub device delete version: %s", devices.Version)
		deviceList := devices.Devices
		for _, deviceInfo := range deviceList {
			glog.Infof("delete sub device. deviceId : %s", deviceInfo.DeviceId)
		}
	}
	// Here is a demonstration of the gateway actively deleting sub-device requests
	deviceIds := make([]string, 1)
	deviceIds = append(deviceIds, "xxxxx")
	gatewayDevice.DeleteSubDevices(deviceIds)
```
Users can implement the SubDevicesDeleteHandler method by themselves. This example provides a default implementation. After the gateway sends a subdevice deletion request to the platform, the platform will notify the gateway to delete the subdevice after deleting the subdevice.

## 4.12 Report device log information
In /samples/log/log_samples.go, it is demonstrated that the device reports log information.
```go
    < create device code here ... >

    for i := 0; i < 10; i++ {
		entry := model.DeviceLogEntry{
			Type: "DEVICE_MESSAGE",
			//Timestamp: iot.GetEventTimeStamp(),
			Content: "message hello " + strconv.Itoa(i),
		}
		entries = append(entries, entry)
	}

	for i := 0; i < 100; i++ {
		result := device.ReportLogs(entries)
		fmt.Println(result)

	}
```
Users can report different log information in different periods of the device according to their needs. For example, a log will be reported after the device is disconnected and reconnected.

## 4.13 End-side rules
Create a client-side rule in the console:
![](.\doc\figure_en\device_rule_en.png)

The ruleManage() method in /samples/rule/rule_demo.go implements an example of end-side rules. You can enable end-side rules through the following parameters
```go
authConfig.RuleEnable = true
```
Set the device command listener by implementing the CommandHandler method, that is:

```go
device.Client.CommandHandler = func(command model.Command) (bool, interface{}) {
		glog.Infof("command device id is %s", command.ObjectDeviceId)
		glog.Infof("command name is %s", command.CommandName)
		glog.Infof("command serviceId is %s", command.ServiceId)
		glog.Infof("command params is %v", command.Paras)
		return true, map[string]interface{}{
			"cost_time": 12,
		}
	}
```

When the rule is triggered, the device will automatically call the method in the listener after receiving the command.
The example prints the content of the command in the CommandHandler method. Customers can customize the content to implement a series of operations on the device.

```python
def run():
    
    < create device code here ... >
    
    # Set up listener
    device.Client.CommandHandler = func(command model.Command) (bool, interface{}) {
		glog.Infof("command device id is %s", command.ObjectDeviceId)
		glog.Infof("command name is %s", command.CommandName)
		glog.Infof("command serviceId is %s", command.ServiceId)
		glog.Infof("command params is %v", command.Paras)
		return true, map[string]interface{}{
			"cost_time": 12,
		}
	}
    
    connect := device.Connect()
	glog.Infof("connect result : %v", connect)

    // Report the SDK version. The client-side rules need to have the SDK version number before they can be created.
	device.ReportDeviceInfo("", "")

    logger.info("begin to report properties")
    //Report attributes
	props := model.DevicePropertyEntry{
		ServiceId: "smokeDetector",
		EventTime: iot.GetEventTimeStamp(),
		Properties: test_model.DemoProperties{
			Temperature: 1,
		},
	}

	var content []model.DevicePropertyEntry
	content = append(content, props)
	services := model.DeviceProperties{
		Services: content,
	}
	device.ReportProperties(services)
	for {
		time.Sleep(1 * time.Second)
	}
```

Executing the run function, the code will produce the following output:
![](.\doc\figure_en\device_rule_action_en.png)

If you want to use a custom method to process the actions of end-side rules, customRuleManage() in /samples/rule/rule_demo.go implements an example of customizing end-side rules.
The RuleActionHandler in the following code is a custom end-side rule processing method, and the instance of RuleActionHandler is set as the command listener, that is:

```go
device.Client.RuleActionHandler = func(actions []model.Action) bool {
		for _, action := range actions {
			glog.Infof("action deviceId: %s:", action.DeviceId)
			command := action.Command
			glog.Infof("action command name : %s", command.CommandName)
			glog.Infof("action command body : %v", command.CommandBody)
		}
		return true
	}
```

When the rule is triggered, the device will automatically call the RuleActionHandler method in the listener when it receives the command.
The example implements some custom operations in the RuleActionHandler method. For example the following output:
![](.\doc\figure_en\device_rule_action_custom_en.png)

## 4.14 Equipment issuance
Create a distribution policy in the console with the keyword xxx:
![](.\doc\figure_en\bootstrap_policy_static_en.png)

Create a device in the console, select static policy as the policy type, and select the product in the same area as the created policy:
![](.\doc\figure_en\bootstrap_create_device_en.png)

The bootstrapSecret() method in /samples/bs/bootstrap_sample.go implements an example of using static policies to allocate devices. Change Servers to the issued address, select the device ID and password just created, and select the device and password.
The value xxx of BaseStrategyKeyword in BootStrapBody is the keyword of the created static strategy. After filling in the correct certificate information.

```go
    //Issue the device ID registered on the platform
	deviceId := "your device id"
	//Device secret key
	pwd := "your device password"

	authConfig := config2.ConnectAuthConfig{
		Id:             deviceId,
		Password:       pwd,
		Servers:        "mqtts://{bootstrap access ip}:8883",
		UseBootstrap:   true,
		BsServerCaPath: "./resource/root.pem",
		ServerCaPath:   "./resource/root.pem",
		BootStrapBody: &amp;model.BootStrapProperties{
			BaseStrategyKeyword: "xxx",
		},
	}
	device := device2.NewMqttDevice(&amp;authConfig)
	if device == nil {
		glog.Warningf("create device failed.")
		return
	}
	initRes := device.Connect()
	glog.Infof("connect result : %v", initRes)
	time.Sleep(3 * time.Second)
	//Report attributes
	device.ReportProperties(test_util.GeneratePropertiesMessage(28))
```
Executing the run function, the code will produce the following output:
![](.\doc\figure_en\bootstrap_device_success_en.png)

In /samples/bs/bootstrap_sample.go, in addition to the bootstrapSecret() method for issuing key authentication devices, there are also the methods bootstrapCert() for issuing certificate devices, the method bootstrapScopeIdSecretStaticPolicy() for issuing key authentication device groups, and the use of certificates. The method bootstrapScopeIdCertStaticPolicy() is used to issue authentication registration group devices.After the device is successfully issued, the access address will be stored in the local server_info.txt file. Next time it is executed, the address in the file will be used first for access. If the device key is updated or other devices are used for provisioning, you need to delete the address. file and then redistribute it.Use static policies for registration group device provisioning. For detailed parameters, please refer to the following link:
[Device access provisioning example](https://support.huaweicloud.com/qs-iotps/iot_03_0006.html)

## 4.15 Disconnection and reconnection
In the connectWithRetry() method in /samples/connect/connect_demo.go, we demonstrate the function of disconnection and reconnection.
```go
    //Close broken link reconnection
    var autoReconnect = false
	authConfig := &amp;config.ConnectAuthConfig{
		Id:       "{your device id}",
		Servers:  "mqtts://{access_address}:8883",
		AuthType:        constants.AuthTypeX509,
        AutoReconnect:   &amp;autoReconnect,
		ServerCaPath: "./resources/root.pem",
        CertFilePath: "your device cert path",
		CertKeyFilePath: "your device cert key path",
	}
	mqttDevice := device.NewMqttDevice(authConfig)
	if mqttDevice == nil {
		glog.Warningf("create mqtt device failed.")
		return
	}
	// After turning off automatic reconnection, you can implement a custom disconnection and reconnection function in the callback function here.
	mqttDevice.Client.ConnectionLostHandler = func(client mqtt.Client, reason error) {
		glog.Warningf("connect lost from server. you can customize auto reconnect logic here")
	}
	// Some custom logic after establishing a link with the platform can be implemented here
	mqttDevice.Client.ConnectHandler = func(client mqtt.Client) {
		glog.Infof("connect to server success.")
	}
	initResult := mqttDevice.Connect()
	glog.Info("connect result is : ", initResult)
	//Report attributes
	mqttDevice.ReportProperties(test_util.GeneratePropertiesMessage(31))
```
You can enable disconnection reconnection by setting config.ConnectAuthConfig.AutoReconnect=True. After setting it to True, you can configure the disconnection reconnection interval and the maximum backoff time by setting the following parameters. The specific logic of disconnection and reconnection can be viewed in the Connect() method in the function in /iot/client/mqtt_device_client.go.
```go
authConfig.MaxBackOffTime = 1000
authConfig.MinBackOffTime = 30 * 1000
authConfig.BackOffTime = 1000
```
You can also set it to False to disable reconnection. Then implement your own disconnection and reconnection logic in a customized way. SDK provides you with the ConnectHandler callback. You can implement this callback function. SDK will notify you when the link is successfully established and the connection is disconnected. You can implement your own link disconnection and reconnection logic in the function and use the following method, Prefab the interface into the sdk. :
```go
// After turning off automatic reconnection, you can implement a custom disconnection and reconnection function in the callback function here.
	mqttDevice.Client.ConnectionLostHandler = func(client mqtt.Client, reason error) {
		glog.Warningf("connect lost from server. you can customize auto reconnect logic here")
	}
```
The sdk also provides the MaxBufferMessage parameter. If you set this parameter, when the sdk is disconnected from the platform, the messages you report will be cached in the memory. The maximum number of cached messages is the value of MaxBufferMessage. If the cached message exceeds this value, The data that enters the cache earliest will be removed. When the link is established with the platform again, the SDK will re-publish the messages in the cache to the platform.
```go
authConfig.MaxBufferMessage = 100
```

#5.0 Open Source License
- Follow the BSD-3 open source license agreement

# 6.0 interface documentation
Refer to [Device Access Interface Documentation](./IoT-Device-SDK-Python-API Documentation.pdf)

# 7.0 More documentation
Refer to [Device Access More Documentation](https://support.huaweicloud.com/devg-iothub/iot_02_0178.html)

