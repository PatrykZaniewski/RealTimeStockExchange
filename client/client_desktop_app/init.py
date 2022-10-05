import asyncio
import websockets

async def hello():
    async with websockets.connect("ws://localhost:5014/ws") as websocket:
        await websocket.send("Hello world!222")
        while True:
            message = await websocket.recv()
            print(message)

if __name__ == "__main__":
    asyncio.run(hello())
