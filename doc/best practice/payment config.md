# Domain
UAT: https://www.paymentstg.horizonpay.cn
Prod: https://www.payment.horizonpay.cn

# 1. login
POST url-> {domain}/payment/account/login
```json
{
    "app_id": "123456",
    "app_secret": "123456",
    "email": "fang.qiang@bindo.com",
    "pwd": "123456"
}
```
```json
{
    "code": "00",
    "data": {
        "email": "fang.qiang@bindo.com",
        "name": "Fang",
        "role": "admin",
        "token": "8aaee146-96f4-4af4-a914-289d8b78d036",
        "user_id": 62
    }
}
```
# 2. Query merchant information：（choose one merchant id to complete the flow process）
POST url-> {domain}/payment/merchant/query
Header: X-Access-Token: "8aaee146-96f4-4af4-a914-289d8b78d036"

```json
{
    "limit": 100,
    "offset": 0
}
```
```json
{
    "code": "00",
    "data": {
        "data": [
            {
                "addr": "123124",
                "agency_id": 6,
                "created_at": "2021-01-23T07:42:02.7818Z",
                "email": "agencytest@163.com",
                "id": 10,
                "name": "Merchant 1",
                "tel": "12312412",
                "updated_at": "2021-01-23T07:42:02.7818Z"
            },
            {
                "agency_id": 7,
                "created_at": "2020-11-27T07:20:08.93688Z",
                "email": "123124@123.com",
                "id": 9,
                "name": "payment",
                "tel": "12312123",
                "updated_at": "2020-11-27T07:20:08.93688Z"
            },
            {
                "addr": "adfadf",
                "agency_id": 4,
                "created_at": "2020-08-27T07:07:07.221136Z",
                "email": "merchant@163.com",
                "id": 8,
                "name": "TestMerchant1",
                "tel": "1231241234",
                "updated_at": "2020-08-27T07:07:07.221136Z"
            },
            {
                "addr": "sdfsd",
                "agency_id": 3,
                "created_at": "2020-08-25T10:03:15.001955Z",
                "email": "111@qq.com",
                "id": 7,
                "name": "Naturel lee",
                "tel": "1212121212",
                "updated_at": "2020-08-25T10:03:21.995837Z"
            },
            {
                "addr": "add",
                "agency_id": 3,
                "created_at": "2020-08-24T07:02:09.054965Z",
                "email": "111@qq.com",
                "id": 6,
                "name": "merchantInAgency",
                "tel": "1212121212",
                "updated_at": "2020-08-24T07:07:47.637106Z"
            },
            {
                "addr": "address",
                "agency_id": 1,
                "created_at": "2020-08-24T06:27:49.045684Z",
                "email": "1234552784@qq.com",
                "id": 1,
                "name": "merch",
                "tel": "123456789",
                "updated_at": "2020-08-24T06:27:49.045684Z"
            },
            {
                "addr": "address",
                "agency_id": 3,
                "created_at": "2020-08-24T06:27:49.045684Z",
                "email": "1234552784@qq.com",
                "id": 2,
                "name": "merch",
                "tel": "123456789",
                "updated_at": "2020-08-24T06:27:49.045684Z"
            },
            {
                "addr": "address",
                "agency_id": 1,
                "created_at": "2020-08-24T06:27:49.045684Z",
                "email": "1234552784@qq.com",
                "id": 3,
                "name": "merch",
                "tel": "123456789",
                "updated_at": "2020-08-24T06:27:49.045684Z"
            },
            {
                "addr": "address",
                "agency_id": 1,
                "created_at": "2020-08-24T06:27:49.045684Z",
                "email": "1234552784@qq.com",
                "id": 4,
                "name": "merch",
                "tel": "123456789",
                "updated_at": "2020-08-24T06:27:49.045684Z"
            },
            {
                "addr": "address",
                "agency_id": 1,
                "created_at": "2020-08-24T06:27:49.045684Z",
                "email": "1234552784@qq.com",
                "id": 5,
                "name": "merch",
                "tel": "123456789",
                "updated_at": "2020-08-24T06:27:49.045684Z"
            }
        ],
        "total": 10
    }
}
```

# 3. 获取支付配置
POST url-> {domain}/payment/config
Header: X-Access-Token: "8aaee146-96f4-4af4-a914-289d8b78d036"
```json
{
    "device_sn": "555555",
    "limit": 0,
    "merchant_id": 9,
    "offset": 0
}
```
```json
{
    "code": "00",
    "data": {
        "data": [
            {
                "acquirer_id": 10,
                "addition": "12312wefaf",
                "created_at": "2021-01-15T11:39:19.648504Z",
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
                "updated_at": "2021-01-15T11:43:46.684804Z"
            },
            {
                "acquirer_id": 10,
                "addition": "asdfawefwf",
                "created_at": "2021-01-15T11:34:19.010968Z",
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
                "updated_at": "2021-01-15T11:34:19.010968Z"
            },
            {
                "created_at": "2021-01-13T14:08:07.707773Z",
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
                "updated_at": "2021-01-13T14:08:07.707773Z"
            }
        ],
        "total": 3
    }
}
```

# 4.获取收单配置：（根据需要循环获取）
POST url-> {domain}/payment/agency_acquirer/query
Header: X-Access-Token: "8aaee146-96f4-4af4-a914-289d8b78d036"
```json
{
    "filters": {
        "id": "10"
    },
    "limit": 100,
    "offset": 0
}
```
```json
{
    "code": "00",
    "data": {
        "data": [
            {
                "addition": "{\n    \"grpc_connect_info\":\"localhost:50002\",\n    \"TPDU\":\"6000170000\",\n    \"NII\":\"017\",\n    \"URL\":\"202.127.169.216:6868\"\n}",
                "agency_id": 0,
                "auto_settlement_time": "",
                "config_file_url": "",
                "created_at": "2020-12-23T11:30:29Z",
                "id": 10,
                "impl_name": "with_tid_default",
                "name": "boc",
                "updated_at": "2020-12-23T11:30:31Z"
            }
        ],
        "total": 1
    }
}
```


