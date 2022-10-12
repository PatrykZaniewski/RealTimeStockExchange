import datetime
import os

import requests as requests

from client.client_desktop_app.model.order import OrderType, Order, OrderSubtype


def process_order(asset_name: str, amount: int, price: float, order_type: OrderType):
    order = Order.create(asset_name, amount, price, order_type.value, OrderSubtype.MARKET_ORDER.value)
    print(f"{order.id},SEND,{datetime.datetime.timestamp(datetime.datetime.now())}")
    requests.post(url=f"{os.getenv('broker_facade_url')}/order", data=order.to_json(),
                  headers={"id": os.getenv("identifier")})
