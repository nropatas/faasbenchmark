function cpuIntensiveCalculation(baseNumber) {
    const iterationCount = 50000 * Math.pow(baseNumber, 3);
    let result = 0;
    for (let i = iterationCount; i >= 0; i--) {
        result += Math.atan(i) * Math.tan(i);
    }
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

exports.cpustress = (args) => {
    const startTime = process.hrtime();
    const params = getParameters(args);
    if (params.error) {
        return {
            statusCode: 200,
            headers: { 'Content-Type': 'application/json' },
            body: { error: params.error },
        };
    }

    runTest(params);

    const reused = isWarm();
    const duration = getDuration(startTime);

    return {
        statusCode: 200,
        headers: {
            'Content-Type': 'application/json',
        },
        body: { reused, duration },
    }
}
