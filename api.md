# api文档

## 基础协议
Header:
|    字段名    | 说明  |
| ----------  | ---  |
|   X-ACCESS-TOKEN   |  login返回回来的token，除了登录api，其他都需要带 |

Response:
```json
{
    "code":0,
    "msg":"OK"
}
```

## 账号部分
### 登录
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

### 登出
***POST***
```https://{base_url}/payment/account/logout```

### 验证登录状态
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

### 添加账号(只有admin才有此权限)
***POST***
```https://{base_url}/payment/account/create```
Request Example:
```json
{
    "email":"123455@qq.com",
    "pwd":"123456",
    "role": "admin",  // "admin" "machine" "normal"
    "name":"xxx"
}
```

### 删除账号(只有admin才有此权限)
***POST***
```https://{base_url}/payment/account/delete```
Request Example:
```json
{
    "id":123
}
```

### 更新账号数据（只发送需要更新的字段）
***POST***
```https://{base_url}/payment/updateaccount```
Request Example:
```json
{
    "id":123,
    "account_num":"123455@qq.com",
    "pwd":"123456",
    "role": "admin",  // "admin" "machine" "normal"    只有管理员才可以更改这个字段
    "name":"xxx",
    "email":"xxx"
}
```

### 批量查询的账号信息
***POST***
```https://{base_url}/payment/queryaccount```
Request Example:
```json
{
    "offset": 1234,
    "limit": 123,
    "filters": {
      "account_num": "123455@qq.com",
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
        "role": "admin",  // "admin" "machine" "normal"    只有管理员才可以更改这个字段
        "name":"xxx",
        "email":"xxx"
      }     
    ]
}
```

## 商户部分
### 获取此账号关联的商户信息
***POST***
```https://{base_url}/payment/merchants```
Response Example:
```json
{
  "merchants": [
    {
      "id": 123,
      "name": "213",
      "tel": "123",
      "addr": "",
      "role": "sales"  // "admin", "sales"，此账号在这个店里面的角色
    }
  ]
}
```
### admin获取关联的商户信息
```https://{base_url}/payment/querymerchants```
Request Example:
```json
{
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
    "merchants": [
      {
            "id": 123,
            "name": "213",
            "tel": "123",
            "addr": ""
      }    
    ]
}
```

### 创建店铺(只有超级管理员才能做)
***POST***
```https://{base_url}/payment/createmerchant```
Request Example:
```json
{
    "name": "213",
    "tel": "123",
    "addr": ""
}
```
Response Example:
```json
{
    "total": 1234,   // 一共有多少店铺
    "merchants":{
            "id": 123,
            "name": "213",
            "tel": "123",
            "addr": ""
      }    
}
```
### 删除店铺(不允许使用)
***POST***
```https://{base_url}/payment/delmerchant```
Request Example:
```json
{
    "id": 123
}
```

### 修改店铺信息
***POST***
```https://{base_url}/payment/updatemerchant```
Request Example:
```json
{
    "id": 123,
    "name": "213",
    "tel": "123",
    "addr": ""  
}
```

### 给店铺添加关联账号(只有超级管理员和店铺管理员才可以使用)
***POST***
```https://{base_url}/payment/addaccount2merchant```
Request Example:
```json
{
    "merchant_id": 123,
    "account_id": 213
}
```

### 查询店铺所有的关联账号
***POST***
```https://{base_url}/payment/account2merchant```
Request Example:
```json
{
    "merchant_id": 123
}
```
Response Example:
```json
{
  "accounts": [
      {
        "id":123,
        "account_num":"123455@qq.com",
        "pwd_hash":"abcd",
        "role": "admin",  // "admin" "machine" "normal"    只有管理员才可以更改这个字段
        "name":"xxx",
        "email":"xxx"
      }     
    ]
}
```


### 给店铺添加关联账号(只有超级管理员和店铺管理员才可以使用)
***POST***
```https://{base_url}/payment/addaccount2merchant```
Request Example:
```json
{
    "merchant_id": 123,
    "account_id": 213
}
```

### 给店铺删除关联账号(只有超级管理员和店铺管理员才可以使用)
***POST***
```https://{base_url}/payment/delaccount2merchant```
Request Example:
```json
{
    "merchant_id": 123,
    "account_id": 213
}
```

### 给店铺添加设备
***POST***
```https://{base_url}/payment/adddevice2merchant```
Request Example:
```json
{
  "merchant_id": 123,
  "device_id": "123"
}
```

## 支付参数部分
### 查询支付参数
***POST***
```https://{base_url}/payment/parameter/query```
Request Example:
```json
{
    "merchant_id": 123,
    "device_id": "123"
}
```

Response Example:
```json
{
  "parameters": [
    {
        "payment_method": ["sale", "void","refund"],
        "payment_type": ["visa", "mastercard", "cup", "amex"],
        "entry_type": ["swipe", "contact", "contactless"],
        "mid":"1234545",
        "tid":"123354",
        "url":"123445",
        "addition":""   // 附加参数
    }
  ]
}
```

### 添加支付参数
***POST***
```https://{base_url}/payment/parameter/add```
Request Example:
```json
{
    "merchant_id": 123,
    "device_id": "123",
    "parameter":{
        "payment_method": ["sale", "void","refund"],
        "payment_type": ["visa", "mastercard", "cup", "amex"],
        "entry_type": ["swipe", "contact", "contactless"],
        "mid":"1234545",
        "tid":"123354",
        "url":"192.168.1.1:8080",
        "addition":""   // 附加参数
    }
}
```

Response Example:
```json
{
  "id": 123
}
```

### 更新支付参数
***POST***
```https://{base_url}/payment/parameter/update```
Request Example:
```json
{
    "id": 123,
    "parameter":{
        "payment_method": ["sale", "void","refund"],
        "payment_type": ["visa", "mastercard", "cup", "amex"],
        "entry_type": ["swipe", "contact", "contactless"],
        "mid":"1234545",
        "tid":"123354",
        "url":"192.168.1.1:8080",
        "addition":""   // 附加参数
    }
}
```

### 删除支付参数
***POST***
```https://{base_url}/payment/parameter/del```
Request Example:
```json
{
    "id": 123
}
```

## 设备管理部分
### 查询设备信息
***POST***
```https://{base_url}/payment/tms/device/query```
Request Example:
```json
{
    "merchant_id": 123,  // 不带这个参数可以查询所有设备（超级管理员才有此权限）

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
        "csn": "123",
        "model": "N5",
        "alias": "123",
        "reboot_model": "every_day",
        "reboot_time": "05:30",
        "reboot_day_in_month": 1,
        "reboot_day_in_week": 1,
        "location_lat": "",
        "location_lin": "",
        "push_token": ""
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
    "model": "N5",
    "alias": "123",
    "reboot_model": "every_day",
    "reboot_time": "05:30",
    "reboot_day_in_month": 1,
    "reboot_day_in_week": 1,
    "location_lat": "",
    "location_lin": "",
    "push_token": ""
}
```

### 查询设备内部app信息
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
        "status": "",
        "app_id": 12,
        "app_file_id": 123
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
    "status": "",
    "app_id": 12,
    "app_file_id": 123
}
```

### 新增设备内部app信息
***POST***
```https://{base_url}/payment/tms/deviceapp/add```
Request Example:
```json
{
    "id": 123,
    "name": "",
    "package_id": "",
    "version_name": "",
    "version_code": "",
    "status": "",
    "app_id": 12,
    "app_file_id": 123
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
```https://{base_url}/payment/tms/app/del```
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
   "devices":[
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
   "devices":[
    {
        "id": 123,
        "file_name": "",
        "version_code": 123,
        "version_name": "v1.0.0",
        "desc": "",
        "decode_status": "",
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





