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
    "email":"123455@qq.com",
    "pwd":"123456",
    "app_id":"123456"
}
```

Response Example:
```json
{
    "token":"123456789",
    "role": "admin",  // "admin" "machine" "user"
    "name":"xxx",
    "email":"xxx"
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
    "role": "admin",  // "admin" "machine" "user"
    "name":"xxx",
    "email":"xxx"
}
```

### 添加账号(只有admin才有此权限)(画面 2)
***POST***
```https://{base_url}/payment/account/create```
Request Example:
```json
{
    "email":"123455@qq.com",
    "pwd":"123456",
    "role": "admin",  // "admin" "machine" "user"
    "name":"xxx"
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
    "offset": 1234,
    "limit": 123,
    "filters": {
      "email": "123455@qq.com",
      "id": 123
    } 
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有多少账号
    "data": [
      {
        "id":123,
        "account_num":"123455@qq.com",
        "pwd_hash":"abcd",
        "role": "admin",  // "admin" "machine" "user"    只有管理员才可以更改这个字段
        "name":"xxx",
        "email":"xxx"
      }     
    ]
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
    "id": 123,
    "name": "asdf",
    "tel": "123",
    "addr": "asdf",
    "email": ""
}
```

### 删除机构(画面 Agency管理-Agency 信息)
***POST***
```https://{base_url}/payment/agency/delete```
Request Example:
```json
{
    "id": 123
}
```

### 查询机构信息(画面 Agency管理-Agency list)
***POST***
```https://{base_url}/payment/agency/query```
Request Example:
```json
{
    "offset": 1234,
    "limit": 123,
    "filters": {
      "name": "123455@qq.com",
      "id": 123
    } 
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有多少记录
    "data": [
      {
        "id":123,
        "name":"123123",
        "tel":"123123",
        "addr": "asdfdf",
        "email": "",
        "update_at":"2020-06-13T15:41:16.142489+08:00"
      }     
    ]
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
    "agency_id": 123,
    "offset": 1234,
    "limit": 123,
    "filters": {
    } 
}
```
Response Example:
```json
{
    "total": 1234,
    "data": [
      {
        "id":123,
        "email":"123123",
        "name":"123123",
        "update_at":"2020-06-13T15:41:16.142489+08:00"
      }     
    ]
}
```

### 添加acquirer(画面 新6-商户-acquirer2)
***POST***
```https://{base_url}/payment/agency_acquirer/add```
Request Example:
```json
{
    "name": "",
    "addition": "",
    "config_file_url": ""
}
```
### 更新acquirer
***POST***
```https://{base_url}/payment/agency_acquirer/update```
Request Example:
```json
{
    "id": 123,
    "name": "",
    "addition": "",
    "config_file_url": ""
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
    "offset": 1234,
    "limit": 123,
    "filters": {
    } 
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有多少记录
    "data": [
      {
        "id":123,
        "name": "",
        "addition": "",
        "config_file_url": "",
        "update_at":"2020-06-13T15:41:16.142489+08:00"
      }     
    ]
}
```

### 获取payment methods
***GET***
```https://{base_url}/payment/agency_acquirer/payment_methods```
Response Example:
```json
{
    "data": [
      "visa","mastercard","jcb","cup"   
    ]
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
    "agency_id": 123,
    "name": "213",
    "tel": "123",
    "addr": "",
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

## 商户部分
### 获取商户信息(画面 8商户管理-商户列表)
```https://{base_url}/payment/merchant/query```
Request Example:
```json
{
    "agency_id": 123,   // 如果是获取所有商户信息，和机构无关，则不传
    "offset": 1234,
    "limit": 123,
    "filters": {
      "name": "kfc",
      "id": 123
    } 
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有多少店铺
    "data": [
      {
            "id": 123,
            "name": "213",
            "tel": "123",
            "email": "",
            "addr": "",
            "update_at":"2020-06-13T15:41:16.142489+08:00"
      }    
    ]
}
```

### 给店铺添加关联账号(只有超级管理员，机构管理员，店铺管理员才可以使用)(画面 新4-商户-新增员工4)
***POST***
```https://{base_url}/payment/merchant_associate/add```
Request Example:
```json
{
    "merchant_id": 123,
    "account_id": 213,
    "role": "admin"   // "admin", "stuff"
}
```

### 查询店铺所有的关联账号(画面 15新增商户 设备 员工 - 员工)
***POST***
```https://{base_url}/payment/merchant_associate/query```
Request Example:
```json
{
    "merchant_id": 123,
    "offset": 1234,
    "limit": 123,
    "filters": {
    }
}
```
Response Example:
```json
{ 
    "total": 1234,   // 一共有多少数据
    "data": [
      {
        "id":123,
        "role": "admin",  // "admin", "stuff"    只有管理员才可以更改这个字段
        "name":"xxx",
        "email":"xxx",
        "updated_at":"2020-06-14T14:34:13.058434+08:00"
      }     
    ]
}
```

### 给店铺更新关联账号(只有超级管理员，机构管理员，店铺管理员才可以使用)(画面 新4-商户-新增员工4)
***POST***
```https://{base_url}/payment/merchant_associate/update```
Request Example:
```json
{
    "id": 123,
    "account_id": 213,
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
    "merchant_id": 123,
    "offset": 1234,
    "limit": 123,
    "filters": {
    }
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有多少店铺
    "data": [
      {
        "id":123,
        "terminal_id": "asdfsdf",
        "cdid": "1234",
        "updated_at":"2020-06-14T14:34:13.058434+08:00"
      }     
    ]
}
```

### 给店铺设备添加支付参数(画面 新3-商户-新增设备3)
***POST***
```https://{base_url}/payment/merchant_device_payment/add```
Request Example:
```json
{
  "merchant_device_id": 123,
  "payment_methods": ["visa","mastercard","jcb"],
  "entry_types": ["swipe","contact","contactless"],
  "payment_types": ["sale", "void", "refund"],
  "acquirer_id": 123,
  "mid": "",
  "tid": "",
  "addition": ""
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
  "id": 123,
  "payment_methods": ["visa","mastercard","jcb"],
  "entry_types": ["swipe","contact","contactless"],
  "payment_types": ["sale", "void", "refund"],
  "acquirer_id": 123,
  "mid": "",
  "tid": "",
  "addition": ""
}
```

### 查询店铺设备支付参数
***POST***
```https://{base_url}/payment/merchant_device_payment/query```
Request Example:
```json
{
    "merchant_id": 123, 
    "device_id": 123,
    "offset": 1234,
    "limit": 123,
    "filters": {
    }
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有多少店铺
    "data": [
      {
        "id": 123,
        "payment_methods": ["visa","mastercard","jcb"],
        "entry_types": ["swipe","contact","contactless"],
        "payment_types": ["sale", "void", "refund"],
        "acquirer_id": 123,
        "mid": "",
        "tid": "",
        "addition": "",
        "updated_at":"2020-06-14T14:34:13.058434+08:00"
      }
    ]
}
```

## 设备管理部分
### 查询设备信息(画面 17设备列表)
***POST***
```https://{base_url}/payment/tms/device/query```
Request Example:
```json
{
    "offset": 1234,
    "limit": 123,
    "filters": {
      "device_id": "1233435"
    } 
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有设备
    "devices":[
      {
        "id": "123",
        "device_sn": "1234134",
        "battery": "123",
        "device_csn": "123",
        "device_model": "N5",
        "alias": "123",
        "reboot_mode": "every_day",
        "reboot_time": "05:30",
        "reboot_day_in_month": 1,
        "reboot_day_in_week": 1,
        "location_lat": "",
        "location_lon": "",
        "push_token": "",
        "tags": [
          {
            "id": 123,
            "name": "123"
          }
        ],
        "updated_at":"2020-06-14T14:34:13.058434+08:00"
      }
    ]
}
```
### 更新设备信息
***POST***
```https://{base_url}/payment/tms/device/update```
Request Example:
```json
{
    "id": "123",
    "csn": "123",
    "alias": "123",
    "reboot_model": "every_day",
    "reboot_time": "05:30",
    "reboot_day_in_month": 1,
    "reboot_day_in_week": 1,
    "tags": [
      {
        "id": 123,
        "name": "123"
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
    "apps":[
      {
        "id": "123",
        "name": "",
        "package_id": "",
        "version_name": "",
        "version_code": "",
        "status": "Pending Install",   // "Pending Install", "Installed", "Pending Uninstall", "Warning Installed"
        "app": {
          "id": 123,
          "package_name": ""
        },
        "app_file": {
          "id": 123,
          "name": ""
        },
        "updated_at":"2020-06-14T14:34:13.058434+08:00"
      }
    ]
}
```

### 更新设备内部app信息
***POST***
```https://{base_url}/payment/tms/deviceapp/update```
Request Example:
```json
{
    "id": 123,
    "name": "",
    "package_id": "",
    "version_name": "",
    "version_code": "",
    "status": "Pending Install",   // "Pending Install", "Installed", "Pending Uninstall", "Warning Installed"
    "app": {
      "id": 123
    },
    "app_file": {
      "id": 123
    }
}
```

### 新增设备内部app信息
***POST***
```https://{base_url}/payment/tms/deviceapp/add```
Request Example:
```json
{
    "name": "",
    "package_id": "",
    "version_name": "",
    "version_code": "",
    "status": "",
    "app": {
      "id": 123
    },
    "app_file": {
      "id": 123
    }
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
    "name": "",
    "package_id": "",
    "desc": ""
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
    "id": 123,
    "name": "",
    "package_id": "",
    "desc": ""
}
```

### 查询app
***POST***
```https://{base_url}/payment/tms/app/query```
Request Example:
```json
{
    "offset": 1234,
    "limit": 123,
    "filters": {
      "id": 123,
      "package_id": ""
    } 
}
```
Response Example:
```json
{
  "total": 1234,   // 一共有app
   "data":[
    {
        "id": 123,
        "name": "",
        "package_id": "",
        "desc": ""
    }
  ]
}
```

### 新增app file
***POST***
```https://{base_url}/payment/tms/appfile/add```
Request Example:
```json
{
    "app_id": "",  // 所属的app id
    "url": "",     // 文件路径
    "desc": ""
}
```

### 删除app
***POST***
```https://{base_url}/payment/tms/appfile/del```
Request Example:
```json
{
    "id": 123
}
```

### 更新app
***POST***
```https://{base_url}/payment/tms/appfile/update```
Request Example:
```json
{
    "id": 123,
    "url": "",     // 文件路径
    "desc": ""
}
```

### 查询app
***POST***
```https://{base_url}/payment/tms/appfile/query```
Request Example:
```json
{
    "app_id": 123,
    "offset": 1234,
    "limit": 123,
    "filters": {
      "id": "",
      "file_name": ""
    } 
}
```
Response Example:
```json
{
  "total": 1234,   // 一共有app file
   "data":[
    {
        "id": 123,
        "file_name": "",
        "version_code": 123,
        "version_name": "v1.0.0",
        "desc": "",
        "decode_status": "",  // "Success", "Fail"
        "decode_error_msg": ""
    }
  ]
}
```

### 新增tag
***POST***
```https://{base_url}/payment/tms/tag/add```
Request Example:
```json
{
    "name": ""
}
```

### 删除tag
***POST***
```https://{base_url}/payment/tms/tag/del```
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
    "name": ""
}
```

### 查询tag
***POST***
```https://{base_url}/payment/tms/tag/query```
Request Example:
```json
{
    "app_id": 123,
    "offset": 1234,
    "limit": 123,
    "filters": {
      "id": "",
      "name": ""
    } 
}
```
Response Example:
```json
{
  "total": 1234,   // 一共有tag
   "devices":[
    {
        "id": 123,
        "name": ""
    }
  ]
}
```

### 新增批量更新任务
***POST***
```https://{base_url}/payment/tms/batchupdata/add```
Request Example:
```json
{
    "tags": ["",""],
    "models": ["",""],
    "desc": "",
    "apps": [
      {
        "status": "",
        "app_id": 12,
        "app_file_id": 123
      }   
    ]
}
```

### 删除batchupdata
***POST***
```https://{base_url}/payment/tms/batchupdata/del```
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
      "id": ""
    } 
}
```
Response Example:
```json
{
  "total": 1234,   // 一共多少记录
   "records":[
    {
        "tags": ["",""],
        "models": ["",""],
        "desc": "",
        "status": "",
        "error_msg": "",
        "apps": [
          {
            "status": "",
            "app_id": 12,
            "app_file_id": 123
          }
        ]
    }
  ]
}
```

## 上送文件管理
### 新建上传文件任务
***POST***
```https://{base_url}/payment/file/add```
Request Example:
```json
{
    "md5": "",
    "file_name": "",
    "file_size": 123,
    "tag": "apk files"
}
```
Response Example:
```json
{
    "upload_url": "https://asdfadfdf",
    "download_url": "https://asdfasdf",
    "expired_at": "2020-06-13T15:41:16.142489+08:00"
}
```
