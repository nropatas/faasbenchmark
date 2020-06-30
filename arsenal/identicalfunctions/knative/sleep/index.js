const bodyParser = require('body-parser')
const express = require('express');

const app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: true,
}));

const wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
    const sleep_time = event.query.sleep ? parseInt(event.query.sleep) : null;
    if (!sleep_time && sleep_time !== 0) {
        return {"error": "invalid sleep parameter"};
    }
    return sleep_time;
}

function getParameters(event) {
    return getSleep(event);
}

async function runTest(sleep_time){
    await wait(sleep_time);
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

app.post('/', async (req, res) => {
    const startTime = process.hrtime();
    const params = getParameters(req);
    if (params.error) {
        res.json({
            error: params.error,
        });
        return;
    }

    await runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    res.status(200).json({
        reused: reused,
        duration: duration,
    });
});

const port = process.env.PORT || 8080;

app.listen(port, () => {
    console.log('Listening on port', port);
});
