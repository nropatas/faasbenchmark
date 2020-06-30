const bodyParser = require('body-parser');
const express = require('express');

const app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: true,
}));

const MEGABYTE = 1024 * 1024;

function memIntensive(level) {
    const available_memory = parseInt(process.env.FUNCTION_MEMORY_SIZE, 10);
    const amountInMB = available_memory - (available_memory / 10) * (4 - level);
    console.log(amountInMB);
    Buffer.alloc(amountInMB * MEGABYTE, 'a');
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.query.level ? parseInt(event.query.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

function runTest(intensityLevel){
    memIntensive(intensityLevel);
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

app.post('/', (req, res) => {
    const startTime = process.hrtime();
    const params = getParameters(req);
    if (params.error) {
        res.json({
            error: params.error,
        });
        return;
    }

    runTest(params);

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
