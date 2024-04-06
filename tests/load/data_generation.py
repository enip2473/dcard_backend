from datetime import datetime, timedelta
import random
import requests


def generate_random_date(start_date, end_date):
    delta = end_date - start_date
    random_days = random.randrange(delta.days) - 1
    random_hours = random.randrange(24)
    random_minutes = random.randrange(60)
    return start_date + timedelta(days=random_days, hours=random_hours, minutes=random_minutes)

def generate_data(num_records):
    data = []
    start_at = datetime.strptime("2024-01-01", "%Y-%m-%d").strftime("%Y-%m-%dT%H:%M:%S.000Z")
    end_at_min = datetime.strptime("2024-05-01", "%Y-%m-%d")
    end_at_max = datetime.strptime("2025-05-01", "%Y-%m-%d")
    genders = ["M", "F"]
    countries = ["TW", "JP", "US"]
    platforms = ["android", "ios", "web"]

    for _ in range(num_records):
        end_at = generate_random_date(end_at_min, end_at_max).strftime("%Y-%m-%dT%H:%M:%S.000Z")
        conditions = {
            "ageStart": random.randint(1, 100),
            "ageEnd": random.randint(1, 100),
            "gender": random.sample(genders, 1),
            "country": random.sample(countries, random.randint(1, len(countries))),
            "platform": random.sample(platforms, random.randint(1, len(platforms)))
        }
        # Ensure ageStart is less than or equal to ageEnd
        if conditions["ageStart"] > conditions["ageEnd"]:
            conditions["ageStart"], conditions["ageEnd"] = conditions["ageEnd"], conditions["ageStart"]


        if random.random() < 0.5:
            conditions.pop("gender", None)
        if random.random() < 0.5:
            conditions.pop("country", None)
        if random.random() < 0.5:
            conditions.pop("platform", None)
        if random.random() < 0.5:
            conditions.pop("ageStart", None)
        if random.random() < 0.5:
            conditions.pop("ageEnd", None)
        
        ad = {
            "title": f"AD {_}",
            "startAt": start_at,
            "endAt": end_at,
            "conditions": conditions
        }

        data.append(ad)

    return data

ads_data = generate_data(1000)

for ad_data in ads_data:
    response = requests.post("http://127.0.0.1:8080/api/v1/ad", json=ad_data)
    if response.status_code != 201:
        print(f"Failed to post data: {response.json()}")

