import asyncio
import base64
import json
import random
import uuid
from typing import Dict

import aiohttp

ASSETS_BOUNDS = {
    "ASSECO": [100.00, 150.00],
    "COMARCH": [50.00, 100.00],
    "CDPROJECT": [25.00, 50.00]
}

PROJECT_ID = "angelic-bond-366421"
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
        "orderPrice": round(random.uniform(float(lower_bound), float(lower_bound) + 15.00), 2)
        if order_type == "BUY" else round(
            random.uniform(float(upper_bound) - 15.00, float(upper_bound)), 2),
        "clientId": "mock_client",
        "brokerId": "mock_broker",
        "id": str(uuid.uuid4())
    }


async def process_limit_order():
    asyncio.create_task(publish(generate_order("ASSECO", "SELL", "LIMIT_ORDER")))
    asyncio.create_task(publish(generate_order("ASSECO", "BUY", "LIMIT_ORDER")))
    asyncio.create_task(publish(generate_order("COMARCH", "SELL", "LIMIT_ORDER")))
    asyncio.create_task(publish(generate_order("COMARCH", "BUY", "LIMIT_ORDER")))
    asyncio.create_task(publish(generate_order("CDPROJECT", "SELL", "LIMIT_ORDER")))
    asyncio.create_task(publish(generate_order("CDPROJECT", "BUY", "LIMIT_ORDER")))
    print("LIMIT_ORDER")


async def process_market_order():
    asyncio.create_task(publish(generate_order("ASSECO", "SELL", "MARKET_ORDER")))
    asyncio.create_task(publish(generate_order("ASSECO", "BUY", "MARKET_ORDER")))
    asyncio.create_task(publish(generate_order("COMARCH", "SELL", "MARKET_ORDER")))
    asyncio.create_task(publish(generate_order("COMARCH", "BUY", "MARKET_ORDER")))
    asyncio.create_task(publish(generate_order("CDPROJECT", "SELL", "MARKET_ORDER")))
    asyncio.create_task(publish(generate_order("CDPROJECT", "BUY", "MARKET_ORDER")))
    print("MARKET_ORDER")


async def publish(data: Dict):
    url = 'https://order-collector-dfksv3hpea-lm.a.run.app/order'
    data = base64.b64encode(json.dumps(data).encode("ascii")).decode("ascii")
    tmp = {
        "message": {
            "data": data,
            "messageId": "6067377733834993",
            "message_id": "6067377733834993",
            "publishTime": "2022-10-24T19:18:35.84Z",
            "publish_time": "2022-10-24T19:18:35.84Z"
        },
        "subscription": "projects/broker-366421/subscriptions/broker.internal.client_orders.sub"
    }

    async with aiohttp.ClientSession() as session:
        async with session.post(url=url, data=json.dumps(tmp)):
            print("order")
            pass


async def main():
    while True:
        await asyncio.gather(
            asyncio.sleep(2),
            process_market_order(),
            process_limit_order()
        )


if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(main())
