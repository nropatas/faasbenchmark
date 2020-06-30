const url = require('url');

function cpuIntensiveCalculation(baseNumber) {
    const iterationCount = 50000 * Math.pow(baseNumber, 3);
    let result = 0;
    for (let i = iterationCount; i >= 0; i--) {
        result += Math.atan(i) * Math.tan(i);
    }
}

function getQueryParams(context) {
    const urlParts = url.parse(context.request.url, true);
    return urlParts.query;
}

function getDuration(startTime) {
    const end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    const intensityLevel = event.level ? parseInt(event.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

function runTest(intensityLevel){
    cpuIntensiveCalculation(intensityLevel);
}

function isWarm() {
    const is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

module.exports = async (context) => {
    const query = getQueryParams(context);

    const startTime = process.hrtime();
    const params = getParameters(query);
    if (params.error) {
        return {
            status: 200,
            body: { error: params.error },
        };
    }

    runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return {
        status: 200,
        body: { reused, duration },
    };
};
