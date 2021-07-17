# api文档

## 基础协议
Header:
|    字段名    | 说明  |
| ----------  | ---  |
|   X-ACCESS-TOKEN   |  login返回回来的token，除了登录api，其他都需要带 |

Response:
```json
{
    "code":"00",
    "msg":"OK",
    "data":{}
}
```

## 账号部分
### 登录(画面 1)
***POST***
```https://{base_url}/payment/account/login```
Request Example:
```json
{
    "user_id": 123,
	"app_id": "123456",
	"app_secret": "123456",
	"email": "fang.qiang7@bindo.com",
	"pwd": "123456"
}
```

Response Example:
```json
{
    "email": "fang.qiang7@bindo.com",
    "name": "Fang",
    "role": "user",   // "admin" "machine" "user"
    "token": "fa3cc6db-2a68-4dfa-a292-b8cf2f3bdb35"
}
```

### 登出(画面 7)
***POST***
```https://{base_url}/payment/account/logout```

### 验证登录状态:本地保存了token的情况下，再次启动时验证token的有效性并且获取用户信息
***POST***
```https://{base_url}/payment/account/validate```

Response Example:
```json
{
    "user_id": 123,
    "email": "fang.qiang7@bindo.com",
    "name": "Fang",
    "role": "user"
}
```

### 添加账号(只有admin才有此权限)(画面 2)
***POST***
```https://{base_url}/payment/account/create```
Request Example:
```json
{
	"email": "fang.qiang9@bindo.com",
	"name": "Fang",
	"pwd": "123456",
	"role": "user"
}
```

### 删除账号(只有admin才有此权限)(画面 2)
***POST***
```https://{base_url}/payment/account/delete```
Request Example:
```json
{
    "id":123
}
```

### 更新账号数据（只发送需要更新的字段）(画面 5)
***POST***
```https://{base_url}/payment/updateaccount```
Request Example:
```json
{
    "id":123,
    "account_num":"123455@qq.com",
    "pwd":"123456",
    "role": "admin",  // "admin" "machine" "user"    只有管理员才可以更改这个字段
    "name":"xxx",
    "email":"xxx"
}
```

### 批量查询的账号信息(画面 2)
***POST***
```https://{base_url}/payment/queryaccount```
Request Example:
```json
{
	"filters": {
		"email": "fang.qiang"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "active": true,
            "agency_id": 4,
            "created_at": "2020-07-31T23:19:37.304278+08:00",
            "email": "fang.qiang8@bindo.com",
            "id": 10,
            "name": "Fang",
            "role": "user",
            "updated_at": "2020-07-31T23:19:37.304278+08:00"
        }
    ],
    "total": 1
}
```

## 机构部分
### 添加机构(画面 Agency管理-Agency 信息)
***POST***
```https://{base_url}/payment/agency/add```
Request Example:
```json
{
    "name": "asdf",
    "tel": "123",
    "addr": "asdf",
    "email": ""
}
```

### 更新机构(画面 Agency管理-Agency 信息)
***POST***
```https://{base_url}/payment/agency/update```
Request Example:
```json
{
	"addr": "wuxicun2",
	"agency_id": 0,
	"id": 7,
	"name": "merc"
}
```

### 查询机构信息(画面 Agency管理-Agency list)
***POST***
```https://{base_url}/payment/agency/query```
Request Example:
```json
{
	"filters": {
		"name": "mer"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "addr": "wuxicun2",
            "email": "",
            "id": 5,
            "name": "merc",
            "tel": "123456789"
        },
        {
            "addr": "wuxicun2",
            "email": "adjfasdf.com",
            "id": 7,
            "name": "merc",
            "tel": "123456789"
        }
    ],
    "total": 2
}
```

### 添加机构账号关联,UI错了，Agency这里也应该有一个一样的画面，但是不需要Role属性(画面 15新增商户 设备 员工 - 员工)
***POST***
```https://{base_url}/payment/agency_associate/add```
Request Example:
```json
{
    "agency_id": 123,
    "user_id": 456
}
```

### 删除机构账号关联(画面 15新增商户 设备 员工 - 员工)
***POST***
```https://{base_url}/payment/agency_associate/delete```
Request Example:
```json
{
    "id": 1
}
```
### 查询机构账号关联(画面 15新增商户 设备 员工 - 员工)
***POST***
```https://{base_url}/payment/agency_associate/query```
```json
{
	"agency_id": 5,
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "created_at": "2020-07-31T23:10:32.004435+08:00",
            "email": "fang.qiang7@bindo.com",
            "id": 4,
            "name": "Fang",
            "updated_at": "2020-07-31T23:10:32.004435+08:00"
        }
    ],
    "total": 1
}
```

### 添加acquirer(画面 新6-商户-acquirer2)
***POST***
```https://{base_url}/payment/agency_acquirer/add```
Request Example:
```json
{
	"addition": "addtion",
	"agency_id": 4,
	"config_file_url": "https://asdfadf",
	"name": "BOC"
}
```
### 更新acquirer
***POST***
```https://{base_url}/payment/agency_acquirer/update```
Request Example:
```json
{
    "id": 123,
    "addition": "addtion",
    "config_file_url": "https://asdfadf",
    "name": "BOC"
}
```
### 删除acquirer
***POST***
```https://{base_url}/payment/agency_acquirer/delete```
Request Example:
```json
{
    "id": 123
}
```
### 查询acquirer(画面 新5-商户-acquirer)
***POST***
```https://{base_url}/payment/agency_acquirer/query```
Request Example:
```json
{
	"agency_id": 4,
	"filters": {
		"name": "BOC"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "addition": "addtion",
            "agency_id": 4,
            "config_file_url": "https://asdfadf",
            "created_at": "2020-08-01T17:12:15.654777+08:00",
            "id": 1,
            "name": "BOC",
            "updated_at": "2020-08-01T17:12:15.654777+08:00"
        },
        {
            "addition": "addtion",
            "agency_id": 4,
            "config_file_url": "https://asdfadf",
            "created_at": "2020-08-06T17:11:44.228448+08:00",
            "id": 4,
            "name": "BOC",
            "updated_at": "2020-08-06T17:11:44.228448+08:00"
        }
    ],
    "total": 2
}
```

### 添加device到acquirer(不需要update)
***POST***
```https://{base_url}/payment/agency_device/add```
Request Example:
```json
{
	"agency_id": 4,
	"device_id": 3,
	"file_url": ""    // device_id或者file_url存在一种，不可同时存在
}
```
### 删除acquirer
***POST***
```https://{base_url}/payment/agency_device/delete```
Request Example:
```json
{
    "id": 123
}
```
### 查询acquirer(画面 新5-商户-acquirer)
***POST***
```https://{base_url}/payment/agency_device/query```
Request Example:
```json
{
	"agency_id": 4,
	"filters": {
		"device_sn": "PAX-"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "agency_id": 4,
            "alias": "",
            "battery": 19,
            "created_at": "2020-04-02T08:42:53.065968+08:00",
            "device_csn": "456789",
            "device_model": "17",
            "device_sn": "PAX-A920-0821157228",
            "id": 1,
            "location_lat": "22.546992",
            "location_lon": "113.94653",
            "push_token": "18071adc0391497837f",
            "reboot_day_in_month": 0,
            "reboot_day_in_week": 0,
            "reboot_mode": "1",
            "reboot_time": "03:00:",
            "updated_at": "2020-08-04T00:04:17.023+08:00"
        },
        {
            "agency_id": 4,
            "alias": "",
            "battery": 96,
            "created_at": "2020-04-06T02:28:02.58338+08:00",
            "device_csn": "17010237",
            "device_model": "17",
            "device_sn": "PAX-A920-0821251436",
            "id": 1505156075807081500,
            "location_lat": "22.282085",
            "location_lon": "114.162749",
            "push_token": "160a3797c89bf7fa890",
            "reboot_day_in_month": 0,
            "reboot_day_in_week": 0,
            "reboot_mode": "1",
            "reboot_time": "03:00:",
            "updated_at": "2020-08-03T16:43:54.384973+08:00"
        },
        {
            "agency_id": 4,
            "alias": "",
            "battery": 94,
            "created_at": "2020-04-06T02:30:27.038397+08:00",
            "device_csn": "17010223",
            "device_model": "17",
            "device_sn": "PAX-A920-0820340310",
            "id": 1505156075807081500,
            "location_lat": "22.31403",
            "location_lon": "114.166594",
            "push_token": "100d855909dba529e70",
            "reboot_day_in_month": 0,
            "reboot_day_in_week": 0,
            "reboot_mode": "1",
            "reboot_time": "03:00:",
            "updated_at": "2020-08-03T16:43:54.512664+08:00"
        }
    ],
    "total": 3
}
```

### 获取payment methods
***GET***
```https://{base_url}/payment/agency_acquirer/payment_methods```
Response Example:
```json
{
    "data": [
        "Swipe",
        "Contact",
        "Contactless"
    ],
    "total": 3
}
```

### 获取entry types
***GET***
```https://{base_url}/payment/agency_acquirer/entry_types```
Response Example:
```json
{
    "data": [
      "swipe","contact","contactless"
    ]
}
```

### 获取payment types
***GET***
```https://{base_url}/payment/agency_acquirer/payment_types```
Response Example:
```json
{
    "data": [
      "sale","void","refund"
    ]
}
```

### 创建店铺（画面 Agency-new merchant）
***POST***
```https://{base_url}/payment/merchant/add```
Request Example:
```json
{
	"addr": "wuxicun",
	"agency_id": 7,
	"name": "merchant 1",
	"tel": "123456789",
    "email": ""
}
```

### 删除店铺(不允许使用)（画面 Agency-new merchant）
***POST***
```https://{base_url}/payment/merchant/delete```
Request Example:
```json
{
    "id": 123
}
```

### 修改店铺信息（画面 Agency-new merchant）
***POST***
```https://{base_url}/payment/merchant/update```
Request Example:
```json
{
    "id": 123,
    "name": "213",
    "tel": "123",
    "addr": ""  
}
```

### 获取Agency下面的商户列表
```https://{base_url}/payment/merchant_in_agency/query```
Request Example:
```json
{
	"agency_id": 4,
	"filters": {
		"name": "merchant"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "addr": "wuxicun",
            "agency_id": 4,
            "created_at": "2020-06-13T17:41:05.302358+08:00",
            "id": 6,
            "name": "merchant 1",
            "tel": "123456789",
            "updated_at": "2020-06-13T17:41:05.302358+08:00"
        },
        {
            "addr": "wuxicun",
            "agency_id": 4,
            "created_at": "2020-08-02T16:03:29.6609+08:00",
            "id": 8,
            "name": "merchant 1",
            "tel": "123456789",
            "updated_at": "2020-08-02T16:03:29.6609+08:00"
        }
    ],
    "total": 2
}
```

## 商户部分
### 获取商户信息(画面 8商户管理-商户列表)
```https://{base_url}/payment/merchant/query```
Request Example:
```json
{
	"filters": {
		"name": "merchant"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "addr": "wuxicun",
            "agency_id": 4,
            "created_at": "2020-06-13T17:41:05.302358+08:00",
            "id": 6,
            "name": "merchant 1",
            "tel": "123456789",
            "updated_at": "2020-06-13T17:41:05.302358+08:00"
        },
        {
            "addr": "wuxicun",
            "agency_id": 4,
            "created_at": "2020-08-02T16:03:29.6609+08:00",
            "id": 8,
            "name": "merchant 1",
            "tel": "123456789",
            "updated_at": "2020-08-02T16:03:29.6609+08:00"
        }
    ],
    "total": 2
}
```

### 给店铺添加关联账号(只有超级管理员，机构管理员，店铺管理员才可以使用)(画面 新4-商户-新增员工4)
***POST***
```https://{base_url}/payment/merchant_associate/add```
Request Example:
```json
{
	"merchant_id": 8,
	"role": "admin",
	"user_id": 10
}
```

### 查询店铺所有的关联账号(画面 15新增商户 设备 员工 - 员工)
***POST***
```https://{base_url}/payment/merchant_associate/query```
Request Example:
```json
{
	"limit": 100,
	"merchant_id": 8,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "created_at": "2020-08-06T17:59:52.861507+08:00",
            "email": "fang.qiang8@bindo.com",
            "id": 9,
            "name": "Fang",
            "role": "admin",
            "updated_at": "2020-08-06T17:59:52.861507+08:00"
        }
    ],
    "total": 1
}
```

### 给店铺更新关联账号(只有超级管理员，机构管理员，店铺管理员才可以使用)(画面 新4-商户-新增员工4)
***POST***
```https://{base_url}/payment/merchant_associate/update```
Request Example:
```json
{
    "id": 123,
    "role": "admin"   // "admin", "stuff"
}
```

### 给店铺删除关联账号(只有超级管理员和店铺管理员才可以使用)(15新增商户 设备 员工 - 员工)
***POST***
```https://{base_url}/payment/merchant_associate/delete```
Request Example:
```json
{
    "id": 123
}
```

### 给店铺添加设备(新-商户-新增设备1)
***POST***
```https://{base_url}/payment/merchant_device/add```
Request Example:
```json
{
  "merchant_id": 123,
  "device_id": 123
}
```

### 给店铺删除设备
***POST***
```https://{base_url}/payment/merchant_device/delete```
Request Example:
```json
{
  "id": 123
}
```

### 给店铺修改设备
***POST***
```https://{base_url}/payment/merchant_device/update```
Request Example:
```json
{
  "id": 123,
  "device_id": 123
}
```

### 查询店铺设备
***POST***
```https://{base_url}/payment/merchant_device/query```
Request Example:
```json
{
	"limit": 100,
	"merchant_id": 8,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "created_at": "2020-08-02T22:23:18.142583+08:00",
            "device_id": 2,
            "device_sn": "L8196RCA8W2136",
            "id": 2,
            "updated_at": "2020-08-02T22:23:18.142583+08:00"
        },
        {
            "created_at": "2020-08-06T18:03:21.40396+08:00",
            "device_id": 2,
            "device_sn": "L8196RCA8W2136",
            "id": 3,
            "updated_at": "2020-08-06T18:03:21.40396+08:00"
        }
    ],
    "total": 2
}
```

### 给店铺设备添加支付参数(画面 新3-商户-新增设备3)
***POST***
```https://{base_url}/payment/merchant_device_payment/add```
Request Example:
```json
{
	"acquirer_id": 1,
	"addition": "http://test.com",
	"entry_types": [
		"Swipe",
		"Contact",
		"Contactless"
	],
	"merchant_device_id": 2,
	"mid": "123456789012345",
	"payment_methods": [
		"Visa",
		"MasterCard",
		"Unionpay"
	],
	"payment_types": [
		"Sale",
		"Void",
		"Refund"
	],
	"tid": "12345678"
}
```

### 给店铺设备添加删除参数
***POST***
```https://{base_url}/payment/merchant_device_payment/delete```
Request Example:
```json
{
  "id": 123
}
```

### 给店铺设备更新支付参数
***POST***
```https://{base_url}/payment/merchant_device_payment/update```
Request Example:
```json
{
	"entry_types": [
		"Swipe",
		"Contact"
	],
	"id": 1,
	"payment_methods": [],
	"payment_types": [
		"Sale",
		"Refund"
	]
}
```

### 查询店铺设备支付参数
***POST***
```https://{base_url}/payment/merchant_device_payment/query```
Request Example:
```json
{
	"device_id": 2,
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "total": 1,   
    "data": [
      {
        "acquirer_id": 1,
        "addition": "http://test.com",
        "created_at": "2020-08-02T23:29:32.367686+08:00",
        "entry_types": [
            "Swipe",
            "Contact"
        ],
        "id": 1,
        "merchant_device_id": 2,
        "mid": "123456789012345",
        "payment_methods": [],
        "payment_types": [
            "Sale",
            "Refund"
        ],
        "tid": "12345678",
        "updated_at": "2020-08-06T18:10:43.875195+08:00"
      }
    ]
}
```

## 设备管理部分（不需要新增设备）
### 查询设备信息(画面 17设备列表)
***POST***
```https://{base_url}/payment/tms/device/query```
Request Example:
```json
{
	"filters": {
		"device_sn": "PAX-A920-0821157228"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "agency_id": 4,
            "alias": "",
            "battery": 19,
            "created_at": "2020-04-02T08:42:53.065968+08:00",
            "device_csn": "456789",
            "device_model": "17",
            "device_sn": "PAX-A920-0821157228",
            "id": 1,
            "location_lat": "22.546992",
            "location_lon": "113.94653",
            "push_token": "18071adc0391497837f",
            "reboot_day_in_month": 0,
            "reboot_day_in_week": 0,
            "reboot_mode": "1",
            "reboot_time": "03:00:",
            "tags": [
                {
                    "agency_id": 748,
                    "created_at": "0001-01-01T00:00:00Z",
                    "id": 25,
                    "name": "tag3",
                    "updated_at": "0001-01-01T00:00:00Z"
                }
            ],
            "updated_at": "2020-08-04T00:04:17.023+08:00"
        }
    ],
    "total": 1
}
```
### 更新设备信息
***POST***
```https://{base_url}/payment/tms/device/update```
Request Example:
```json
{
	"alias": "",
	"device_csn": "456789",
	"id": 1,
	"reboot_day_in_month": 0,
	"reboot_day_in_week": 0,
	"reboot_mode": "",
	"reboot_time": "",
	"tags": [
		{
			"id": 25
		}
	]
}
```

### 查询设备内部app信息(画面 新10新增内部APP信息)
***POST***
```https://{base_url}/payment/tms/deviceapp/query```
Request Example:
```json
{
    "device_id": "123"
}
```
Response Example:
```json
{
    "data": [
        {
            "app": null,
            "app_file": null,
            "app_file_id": 0,
            "app_id": 0,
            "created_at": "2020-08-04T16:50:04.894883+08:00",
            "external_id": 1,
            "external_id_type": "merchantdevice",
            "id": 1505221047908102100,
            "name": "",
            "package_id": "",
            "status": "pending uninstall",
            "updated_at": "2020-08-04T16:53:18.815934+08:00",
            "version_code": 0,
            "version_name": ""
        }
    ],
    "total": 1
}
```

### 更新设备内部app信息
***POST***
```https://{base_url}/payment/tms/deviceapp/update```
Request Example:
```json
{
    "app_file_id": 0,
    "app_id": 0,
    "external_id": 1,   // device ID
    "status": "warning installed"  // pending install/ installed/ pending uninstall/ warning installed
}
```

### 新增设备内部app信息
***POST***
```https://{base_url}/payment/tms/deviceapp/add```
Request Example:
```json
{
	"app_file_id": 0,
	"app_id": 0,
	"external_id": 1,   // device ID
	"status": "warning installed"  // pending install/ installed/ pending uninstall/ warning installed
}
```

### 删除设备内部app
***POST***
```https://{base_url}/payment/tms/deviceapp/delete```
Request Example:
```json
{
    "id": 123
}
```

### 新增app
***POST***
```https://{base_url}/payment/tms/app/add```
Request Example:
```json
{
	"description": "测试",
	"name": "Fang Apk",
	"package_id": "com.bindo.test"
}
```

### 删除app
***POST***
```https://{base_url}/payment/tms/app/delete```
Request Example:
```json
{
    "id": 123
}
```

### 更新app
***POST***
```https://{base_url}/payment/tms/app/update```
Request Example:
```json
{
    "description": "MDM",
    "id": 1,
    "name": "MDM",
    "package_id": "com.bindo.mdm"
}
```

### 查询app
***POST***
```https://{base_url}/payment/tms/app/query```
Request Example:
```json
{
	"limit": 100,
	"offset": 0,
    "filters": {
		"name": "MDM"
	}
}
```
Response Example:
```json
{
    "data": [
        {
            "created_at": "2020-08-04T18:18:15.472834+08:00",
            "description": "MDM",
            "id": 1,
            "name": "MDM",
            "package_id": "com.bindo.mdm",
            "updated_at": "2020-08-06T18:23:51.266732+08:00"
        },
        {
            "created_at": "2020-08-06T18:23:11.807733+08:00",
            "description": "测试",
            "id": 1505150659148678100,
            "name": "Fang Apk",
            "package_id": "com.bindo.test",
            "updated_at": "2020-08-06T18:23:11.807733+08:00"
        }
    ],
    "total": 2
}
```

### 新增app file
***POST***
```https://{base_url}/payment/tms/appfile/add```
Request Example:
```json
{
	"app_id": 1,
	"file_url": "https://mdmfiles.oss-cn-hongkong.aliyuncs.com/other%20file/Landi-MDM-V1.15_alpha_release_20200720%20%281%29.apk",
	"update_description": "MDM First Time"
}
```

### 删除app file
***POST***
```https://{base_url}/payment/tms/appfile/delete```
Request Example:
```json
{
    "id": 123
}
```

### 更新app file
***POST***
```https://{base_url}/payment/tms/appfile/update```
Request Example:
```json
{
    "id": 1,
    "update_description": "MDM First Time"
}
```

### 查询app file
***POST***
```https://{base_url}/payment/tms/appfile/query```
Request Example:
```json
{
	"app_id": 1,
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "app_id": 1,
            "created_at": "2020-08-04T21:02:19+08:00",
            "decode_fail_msg": "",
            "file_name": "Landi-MDM-V1.15_alpha_release_20200720%20%281%29.apk",
            "file_url": "https://mdmfiles.oss-cn-hongkong.aliyuncs.com/other%20file/Landi-MDM-V1.15_alpha_release_20200720%20%281%29.apk",
            "id": 1505154471485801700,
            "status": "done",
            "update_description": "MDM First Time",
            "updated_at": "2020-08-04T21:03:43+08:00",
            "version_code": 16,
            "version_name": "MDM-V1.15"
        },
        {
            "app_id": 1,
            "created_at": "2020-08-06T18:27:00+08:00",
            "decode_fail_msg": "",
            "file_name": "Landi-MDM-V1.15_alpha_release_20200720%20%281%29.apk",
            "file_url": "https://mdmfiles.oss-cn-hongkong.aliyuncs.com/other%20file/Landi-MDM-V1.15_alpha_release_20200720%20%281%29.apk",
            "id": 1505154471485801700,
            "status": "done",
            "update_description": "MDM First Time",
            "updated_at": "2020-08-06T18:27:03+08:00",
            "version_code": 16,
            "version_name": "MDM-V1.15"
        }
    ],
    "total": 2
}
```

### 新增tag
***POST***
```https://{base_url}/payment/tms/tag/add```
Request Example:
```json
{
    "name": "Tag2",
    "description":"123123"
}
```

### 删除tag
***POST***
```https://{base_url}/payment/tms/tag/delete```
Request Example:
```json
{
    "id": 123
}
```

### 更新tag
***POST***
```https://{base_url}/payment/tms/tag/update```
Request Example:
```json
{
    "id": 123,
    "name": "",
    "description":"123123"
}
```

### 查询tag
***POST***
```https://{base_url}/payment/tms/tag/query```
Request Example:
```json
{
	"filters": {
		"name": "tag"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "created_at": "0001-01-01T00:00:00Z",
            "id": 23,
            "name": "tag1",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        {
            "created_at": "0001-01-01T00:00:00Z",
            "id": 24,
            "name": "tag2",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        {
            "created_at": "0001-01-01T00:00:00Z",
            "id": 25,
            "name": "tag3",
            "updated_at": "0001-01-01T00:00:00Z"
        }
    ],
    "total": 3
}
```

### 查询tag下的设备
***POST***
```https://{base_url}/payment/tms/device_in_tag/query```
Request Example:
```json
{
	"limit": 100,
	"offset": 0,
    "tag_id": 49
}
```
Response Example:
```json
{
	"code": "00",
	"data": {
		"data": [
			{
				"agency_id": 10,
				"alias": "",
				"battery": 0,
				"created_at": "2021-06-25T11:16:43.326383+08:00",
				"device_csn": "",
				"device_model": "",
				"device_sn": "10202211980114",
				"id": 18902,
				"location_lat": "",
				"location_lon": "",
				"push_token": "",
				"reboot_day_in_month": 0,
				"reboot_day_in_week": 0,
				"reboot_mode": "every_day",
				"reboot_time": "03:00",
				"updated_at": "2021-06-25T11:16:43.326383+08:00"
			}
		],
		"total": 1
	}
}
```

### 新增model
***POST***
```https://{base_url}/payment/tms/model/add```
Request Example:
```json
{
    "name": "Model1"
}
```

### 删除model
***POST***
```https://{base_url}/payment/tms/model/delete```
Request Example:
```json
{
    "id": 123
}
```

### 更新model
***POST***
```https://{base_url}/payment/tms/model/update```
Request Example:
```json
{
    "id": 123,
    "name": ""
}
```

### 查询model
***POST***
```https://{base_url}/payment/tms/model/query```
Request Example:
```json
{
	"filters": {
		"name": "tag"
	},
	"limit": 100,
	"offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "created_at": "0001-01-01T00:00:00Z",
            "id": 23,
            "name": "tag1",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        {
            "created_at": "0001-01-01T00:00:00Z",
            "id": 24,
            "name": "tag2",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        {
            "created_at": "0001-01-01T00:00:00Z",
            "id": 25,
            "name": "tag3",
            "updated_at": "0001-01-01T00:00:00Z"
        }
    ],
    "total": 3
}
```

### 新增批量更新任务
***POST***
```https://{base_url}/payment/tms/batchupdate/add```
Request Example:
```json
{
	"description": "Test1",
	"device_models": [
		{
			"Name": "A920",
			"id": 1
		},
		{
			"Name": "K11",
			"id": 2
		}
	],
	"tags": [
		{
			"id": 1,
			"name": "Tag1"
		},
		{
			"id": 2,
			"name": "Tag2"
		}
	]
}
```

### 删除batchupdata
***POST***
```https://{base_url}/payment/tms/batchupdata/delete```
Request Example:
```json
{
    "id": 123
}
```

### 查询batchupdata
***POST***
```https://{base_url}/payment/tms/batchupdata/query```
Request Example:
```json
{
    "offset": 1234,
    "limit": 123,
    "filters": {
      "name": ""
    } 
}
```
Response Example:
```json
{
    "data": [
        {
            "agency_id": 0,
            "created_at": "2020-08-26T21:53:21.285089+08:00",
            "description": "Test1",
            "device_models": [
                {
                    "Name": "A920",
                    "id": 1
                },
                {
                    "Name": "K11",
                    "id": 2
                }
            ],
            "id": 3,
            "status": "",
            "tags": [
                {
                    "id": 1,
                    "name": "Tag1"
                },
                {
                    "id": 2,
                    "name": "Tag2"
                }
            ],
            "update_fail_msg": "",
            "updated_at": "2020-08-26T21:53:21.285089+08:00"
        }
    ],
    "total": 1
}
```

### 更新批量更新任务
***POST***
```https://{base_url}/payment/tms/batchupdate/update```
Request Example:
```json
{
	"description": "Test2",
	"device_models": [
		{
			"id": 3
		}
	],
	"id": 3,
	"tags": [
		{
			"id": 1
		}
	]
}
```

### 启动批量更新任务
***POST***
```https://{base_url}/payment/tms/batchupdate/starthandle```
Request Example:
```json
{
	"id": 3
}
```

### 查询批量更新内部app信息(功能界面和terminal内部app信息一致)
***POST***
```https://{base_url}/payment/tms/appinbatchupdate/query```
Request Example:
```json
{
    "batch_id": 3,
    "limit": 100,
    "offset": 0
}
```
Response Example:
```json
{
    "data": [
        {
            "app": {
                "agency_id": 0,
                "created_at": "2020-08-24T20:53:05+08:00",
                "description": "",
                "id": 7,
                "name": "Driver",
                "package_id": "com.bindo.uniform.driver",
                "updated_at": "2020-08-24T20:53:11+08:00"
            },
            "app_file": {
                "app_id": 7,
                "created_at": "2020-08-25T04:10:39+08:00",
                "decode_fail_msg": "",
                "file_name": "Landi-uniform-driver-app_v2.0.17_release-20-07-30-18-23.apk",
                "file_url": "https://horizonpay.s3.ap-northeast-2.amazonaws.com/appfile/dbb9e1a7a2ba462f911d564d9661fc60/Landi-uniform-driver-app_v2.0.17_release-20-07-30-18-23.apk",
                "id": 8,
                "status": "done",
                "update_description": "123",
                "updated_at": "2020-08-25T04:14:31+08:00",
                "version_code": 40,
                "version_name": "2.0.17"
            },
            "app_file_id": 8,
            "app_id": 7,
            "created_at": "2020-08-26T22:44:24.516081+08:00",
            "external_id": 3,
            "external_id_type": "batch",
            "id": 44,
            "name": "Driver",
            "package_id": "com.bindo.uniform.driver",
            "status": "pending_install",
            "updated_at": "2020-08-26T22:44:24.516081+08:00",
            "version_code": 40,
            "version_name": "2.0.17"
        }
    ],
    "total": 1
}
```

### 更新批量更新内部app信息
***POST***
```https://{base_url}/payment/tms/appinbatchupdate/update```
Request Example:
```json
{
    "app_file_id": 0,
    "app_id": 0,
    "external_id": 1,   // device ID
    "status": "warning installed"  // pending install/ installed/ pending uninstall/ warning installed
}
```

### 新增批量更新内部app信息
***POST***
```https://{base_url}/payment/tms/appinbatchupdate/add```
Request Example:
```json
{
	"app_file_id": 0,
	"app_id": 0,
	"external_id": 1,   // device ID
	"status": "warning installed"  // pending install/ installed/ pending uninstall/ warning installed
}
```

### 删除批量更新内部app
***POST***
```https://{base_url}/payment/tms/appinbatchupdate/delete```
Request Example:
```json
{
    "id": 123
}
```

### 新增uploadfile（无update操作）
***POST***
```https://{base_url}/payment/tms/uploadfile/add```
Request Example:
```json
{
	"device_sn": "PAX-A920-0821157228",
	"file_name": "12314",
	"file_url": "https://baidu.com/12314"
}
```

### 删除uploadfile
***POST***
```https://{base_url}/payment/tms/uploadfile/delete```
Request Example:
```json
{
    "id": 123
}
```

### 查询uploadfile
***POST***
```https://{base_url}/payment/tms/uploadfile/query```
Request Example:
```json
{
	"limit": 100,
	"offset": 0,
    "filters": {
		"name": "MDM"
	}
}
```
Response Example:
```json
{
    "data": [
        {
            "agency_id": 4,
            "created_at": "2020-08-05T18:39:01.655595+08:00",
            "device_sn": "PAX-A920-0821157228",
            "file_name": "12314",
            "file_url": "https://baidu.com/12314",
            "id": 1,
            "updated_at": "2020-08-05T18:39:01.655595+08:00"
        },
        {
            "agency_id": 4,
            "created_at": "2020-08-06T18:35:19.057972+08:00",
            "device_sn": "PAX-A920-0821157228",
            "file_name": "12314",
            "file_url": "https://baidu.com/12314",
            "id": 734760,
            "updated_at": "2020-08-06T18:35:19.057972+08:00"
        },
        {
            "agency_id": 4,
            "created_at": "2020-08-06T18:47:34.834438+08:00",
            "device_sn": "PAX-A920-0821157228",
            "file_name": "12314",
            "file_url": "https://baidu.com/12314",
            "id": 734761,
            "updated_at": "2020-08-06T18:47:34.834438+08:00"
        }
    ],
    "total": 3
}
```

## 上送文件管理
### 新建上传文件任务
***POST***
```https://{base_url}/payment/file/add```
Request Example:
```json
{
	"file_name": "test1",
	"file_size": 1000000,
	"md5": "1231241234",
	"tag": "appfile"    // app文件上传时，使用appfile；    批量添加设备到agency时：使用 devicefile;   tms里面的文件上送时，使用uploadfile
}
```
Response Example:
```json
{
    "download_url": "https://horizonpay.ap-northeast-2/tms/ad2b123ddd054d9e82111ab9ae05aecf/test1",
    // base64 encoded
    "upload_url": "aHR0cHM6Ly9ob3Jpem9ucGF5LnMzLmFwLW5vcnRoZWFzdC0yLmFtYXpvbmF3cy5jb20vdG1zL2FkMmIxMjNkZGQwNTRkOWU4MjExMWFiOWFlMDVhZWNmL3Rlc3QxP1gtQW16LUFsZ29yaXRobT1BV1M0LUhNQUMtU0hBMjU2JlgtQW16LUNyZWRlbnRpYWw9QUtJQUlSVEpDU1hGRFRPNElGQlElMkYyMDIwMDgwNiUyRmFwLW5vcnRoZWFzdC0yJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDIwMDgwNlQxMDQyMzFaJlgtQW16LUV4cGlyZXM9OTAwJlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCUzQngtYW16LWFjbCZYLUFtei1TaWduYXR1cmU9ZDA0M2M1NDk4MzMwNGM3NDZjNTM5NjQ2MGE5NjcwMzU4MGZiOGFkYmM2OTM2ZDliZTM0YTQzYzRkMWQyYWJhMg=="
    "expired_at": "2020-06-13T15:41:16.142489+08:00"
}
```


## 支付模块
### 查询设备支付参数
***POST***
```https://{base_url}/payment/config```
Request Example:
```json
{
	"merchant_id": 9,
	"device_sn": "555555"
}
```
Response Example:
```json
{
	"code": "00",
	"data": {
		"data": [
			{
				"acquirer_config": {
					"addition": "{\n    \"grpc_connect_info\":\"localhost:50002\",\n    \"TPDU\":\"6000170000\",\n    \"NII\":\"017\",\n    \"URL\":\"202.127.169.216:6868\"\n}",
					"agency_id": 0,
					"auto_settlement_time": "",
					"config_file_url": "",
					"created_at": "2020-12-23T11:30:29+08:00",
					"id": 10,
					"impl_name": "with_tid_default",
					"name": "boc",
					"updated_at": "2020-12-23T11:30:31+08:00"
				},
				"acquirer_id": 10,
				"addition": "12312wefaf",
				"created_at": "2021-01-15T11:39:19.648504+08:00",
				"entry_types": [
					"Swipe",
					"Contact"
				],
				"id": 15,
				"merchant_device_id": 13,
				"mid": "123445",
				"payment_methods": [
					"Visa",
					"MasterCard"
				],
				"payment_types": [
					"Sale"
				],
				"tid": "123123",
				"updated_at": "2021-01-15T11:43:46.684804+08:00"
			},
			{
				"acquirer_config": {
					"addition": "{\n    \"grpc_connect_info\":\"localhost:50002\",\n    \"TPDU\":\"6000170000\",\n    \"NII\":\"017\",\n    \"URL\":\"202.127.169.216:6868\"\n}",
					"agency_id": 0,
					"auto_settlement_time": "",
					"config_file_url": "",
					"created_at": "2020-12-23T11:30:29+08:00",
					"id": 10,
					"impl_name": "with_tid_default",
					"name": "boc",
					"updated_at": "2020-12-23T11:30:31+08:00"
				},
				"acquirer_id": 10,
				"addition": "asdfawefwf",
				"created_at": "2021-01-15T11:34:19.010968+08:00",
				"entry_types": [
					"Swipe"
				],
				"id": 13,
				"merchant_device_id": 13,
				"mid": "123456789",
				"payment_methods": [
					"Visa",
					"MasterCard"
				],
				"payment_types": [
					"Sale"
				],
				"tid": "12345678",
				"updated_at": "2021-01-15T11:34:19.010968+08:00"
			},
			{
				"created_at": "2021-01-13T14:08:07.707773+08:00",
				"entry_types": [
					"Swipe"
				],
				"id": 11,
				"merchant_device_id": 13,
				"mid": "123456789012345",
				"payment_methods": [
					"Visa"
				],
				"payment_types": [
					"Sale"
				],
				"tid": "12345678",
				"updated_at": "2021-01-13T14:08:07.707773+08:00"
			}
		],
		"total": 3
	}
}
```

