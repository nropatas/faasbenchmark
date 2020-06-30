import os
import time
import datetime

from flask import Flask, request, jsonify

app = Flask(__name__)

def get_current_epoch():
    return int((datetime.datetime.utcnow() - datetime.datetime(1970, 1, 1)).total_seconds() * 1000)

def get_sleep_parameter(sleep_param):
    user_input = str(sleep_param)
    if not user_input or not user_input.isdigit() or int(user_input) < 0:
        return {"error": "invalid sleep parameter"}
    return int(user_input)

def run_test(sleep_time):
    time.sleep(sleep_time / 1000.0)

def is_warm():
    is_warm = os.environ.get("warm") == "true"
    os.environ["warm"] = "true"
    return is_warm

@app.route('/')
def sleep():
    start = get_current_epoch()
    reused = is_warm()

    sleep_time = get_sleep_parameter(request.args.get("sleep"))
    if type(sleep_time) != int:
        return jsonify(sleep_time)

    run_test(sleep_time)
    duration = (get_current_epoch() - start) * 1000000

    return jsonify({
        "duration": duration,
        "reused": reused
    })

if __name__ == "__main__":
    app.run(debug=True,host='0.0.0.0',port=int(os.environ.get('PORT', 8080)))