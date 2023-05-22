import os
import time
import pyroscope
from fastapi import FastAPI
from lib.bike.bike import order_bike
from lib.car.car import order_car
from lib.scooter.scooter import order_scooter

application_name = os.getenv("PYROSCOPE_APPLICATION_NAME", "flask-ride-sharing-app")

# Note: If using Grafana Cloud Profiles you'll need to replace the server_address with the one provided in the UI
server_address = os.getenv("PYROSCOPE_SERVER_ADDRESS", "http://pyroscope:4040")
basic_auth_username = os.getenv("PYROSCOPE_BASIC_AUTH_USERNAME", "")
basic_auth_password = os.getenv("PYROSCOPE_BASIC_AUTH_PASSWORD", "")

pyroscope.configure(
	application_name = application_name,
    server_address = server_address,
    basic_auth_username=basic_auth_username, # If using Grafana Cloud profiles
    basic_auth_password=basic_auth_password, # If using Grafana Cloud profiles
	tags             = {
        "region":   f'{os.getenv("REGION")}',
	}
)


app = FastAPI()

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.get("/bike")
def bike():
    order_bike(0.2)
    return "<p>Bike ordered</p>"


@app.get("/scooter")
def scooter():
    order_scooter(0.3)
    return "<p>Scooter ordered</p>"


@app.get("/car")
def car():
    order_car(0.4)
    return "<p>Car ordered</p>"


@app.get("/")
def environment():
    result = "<h1>environment vars:</h1>"
    for key, value in os.environ.items():
        result +=f"<p>{key}={value}</p>"
    return result
