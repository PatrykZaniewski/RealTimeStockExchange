import asyncio
import json
import random
import uuid
from typing import Dict

from google.cloud import pubsub_v1

ASSETS_BOUNDS = {
    "ASSECO": [100.00, 150.00],
    "COMARCH": [50.00, 100.00],
    "CDPROJECT": [25.00, 50.00]
}

PROJECT_ID = "citric-campaign-349210"
TOPIC_ID = "broker_mock.orders"


def generate_order(asset_name: str, order_type: str, order_subtype: str):
    bounds = ASSETS_BOUNDS[asset_name]
    lower_bound = bounds[0]
    upper_bound = bounds[1]
    return {
        "assetName": asset_name,
        "quantity": 1,
        "orderType": order_type,
        "orderSubtype": order_subtype,
        "orderPrice": round(random.uniform(float(upper_bound), float(upper_bound) + 15.00), 2)
        if order_type == "BUY" else round(
            random.uniform(float(lower_bound) - 15.00, float(lower_bound)), 2),
        "clientId": "mock_client",
        "brokerId": "mock_broker",
        "id": str(uuid.uuid4())
    }


async def process_limit_order():
    while True:
        asyncio.create_task(publish(generate_order("ASSECO", "SELL", "LIMIT_ORDER")))
        asyncio.create_task(publish(generate_order("ASSECO", "BUY", "LIMIT_ORDER")))
        await asyncio.sleep(1)


async def process_market_order():
    while True:
        asyncio.create_task(publish(generate_order("ASSECO", "SELL", "MARKET_ORDER")))
        asyncio.create_task(publish(generate_order("ASSECO", "BUY", "MARKET_ORDER")))
        await asyncio.sleep(10)


async def publish(data: Dict):
    publisher = pubsub_v1.PublisherClient()
    topic_path = publisher.topic_path(PROJECT_ID, TOPIC_ID)

    future = publisher.publish(topic_path, json.dumps(data).encode("utf-8"))
    future.result()


async def main():
    await asyncio.gather(process_market_order(), process_limit_order())


if __name__ == "__main__":
    asyncio.run(main())
