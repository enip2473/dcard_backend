# API Documentation

## Create Ad

```
POST /api/v1/ad
```

### Description

Create a new ad.

### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| title | string | The title of the ad.|
| startAt | string | The start date of the ad in ISO 8601 format.|
| endAt | string | The end date of the ad in ISO 8601 format.|
| conditions | object | The conditions of the ad.

### Request

```json
{
    "title": "TW_androi888d",
    "startAt": "2023-12-11T03:00:00.000Z",
    "endAt": "2024-10-27T16:00:00.000Z",
    "conditions": {  
        "ageStart": 20,
        "ageEnd": 30,
        "country": ["TW"],
        "platform": ["android"]
    } 
}
```

### Response

```json
{
    "ID": 1264,
    "CreatedAt": "2024-04-07T06:01:35.514545165Z",
    "UpdatedAt": "2024-04-07T06:01:35.514545165Z",
    "DeletedAt": null,
    "title": "TW_androi888d",
    "startAt": "2023-12-11T03:00:00Z",
    "endAt": "2024-10-27T16:00:00Z",
    "startAge": 20,
    "endAge": 30,
    "genderTarget": "",
    "Countries": [
        {
            "ID": 1,
            "CreatedAt": "2024-04-06T14:59:42.702701Z",
            "UpdatedAt": "2024-04-06T14:59:42.702701Z",
            "DeletedAt": null,
            "code": "TW"
        }
    ],
    "Platforms": [
        {
            "ID": 4,
            "CreatedAt": "2024-04-06T14:59:46.133848Z",
            "UpdatedAt": "2024-04-06T14:59:46.133848Z",
            "DeletedAt": null,
            "name": "android"
        }
    ]
}
```

## List Ads

```
GET /api/v1/ad
```

### Description

List ads with conditions.

### Query Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| offset | int | The offset of the first ad to return, default is 0.|
| limit | int | The maximum number of ads to return (1-100), default is 5.|
| age | int | The age of the ads to return (1-100), default is -1 (all ages).|
| gender | string | The gender (M, F) of the ads to return, default is "" (all genders).|
| country | string | The country (US, CA, JP, TW) of the ads to return, default is "" (all countries).|
| platform | string | The platform (ios, android, web) of the ads to return, default is "" (all platforms).|

### Response

```json
{
    "items": [
        {
            "title": "AD 132",
            "endAt": "2024-04-30T11:33:00+08:00"
        },
        {
            "title": "AD 776",
            "endAt": "2024-04-30T21:00:00+08:00"
        },
        {
            "title": "AD 50",
            "endAt": "2024-05-01T13:31:00+08:00"
        },
        {
            "title": "AD 395",
            "endAt": "2024-05-02T02:59:00+08:00"
        },
        {
            "title": "AD 305",
            "endAt": "2024-05-02T17:49:00+08:00"
        }
    ]
}
```