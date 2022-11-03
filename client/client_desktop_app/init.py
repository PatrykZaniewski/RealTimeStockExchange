import asyncio
import json
import logging
import sys
import time
import uuid

import aiohttp as aiohttp
import requests as requests
import websockets as websockets
from dotenv import load_dotenv
from PyQt5.QtWidgets import QApplication

from client.client_desktop_app.gui import MainWindow
from client.client_desktop_app.websocket_setup import QClient


async def hello():
    async with websockets.connect("ws://broker-366421.ew.r.appspot.com/ws",
                                  extra_headers={"identifier": "broker_client"}) as websocket:
        while True:
            await websocket.send("broker_client")
            message = await websocket.recv()
            message = json.loads(message)
            if message['type'] == 'ORDER_STATUS':
                logging.info(f'{message["id"]},CLIENT,STATUS_RECEIVED,{int(time.time() * 1000000)}')
            # print(message)
            # await websocket.ping("Hello world!")


async def generate_orders():
    # for _ in range(10):
    id = str(uuid.uuid4())
    order = {"assetName": "ASSECO", "quantity": 1, "orderType": "BUY", "orderSubtype": "MARKET_ORDER",
             "orderPrice": 100.0, "id": id}
    tmp = json.dumps(order)
    # logging.info(f'{id},CLIENT,ORDER_SEND,{int(time.time() * 1000000)}')
    async with aiohttp.ClientSession() as session:
        logging.info(f'{id},CLIENT,ORDER_SEND,{int(time.time() * 1000000)}')
        async with session.post(url='https://broker-facade-msdaaqs4fq-lm.a.run.app/order', headers={"identifier": "broker_client"}, data=tmp) as response:
            pass
        # requests.post(url=f"http://localhost:5012/order", data=json.dumps(order), headers={"identifier": "broker_client"})
        # await asyncio.sleep(0.1)


async def main():
    for _ in range(10):
        asyncio.get_event_loop().create_task(generate_orders())
        await asyncio.sleep(1)
    await hello()


def run_gui():
    app = QApplication([])
    app.setQuitOnLastWindowClosed(True)
    window = MainWindow()
    load_dotenv(dotenv_path="./settings.env")
    qclient = QClient(app)
    app.exec_()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO, format='%(asctime)s %(message)s', datefmt='%Y/%m/%d %H:%M:%S')

    run_gui()
    # asyncio.get_event_loop().run_until_complete(main())
