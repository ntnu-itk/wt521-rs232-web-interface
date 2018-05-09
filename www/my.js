// These are set in the HTML template:
//var MAX_LINES = ??
//var dT
var currentState = null

function refresh() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "json?wait");
    xhr.onload = function(response) {
        refreshIntervalId = setTimeout(refresh, 500)

        console.log("Response.target:", response.target)

        if (response.target.status === 200) {
            var state = JSON.parse(response.target.response)
            console.log("Parsed state:", state)
            if (state !== null) {
                currentState = state
            }
        }
        if (currentState !== null) {
            updateBody(currentState)
        }
    }
    xhr.onerror = function(a1, a2) {
        console.log("xhr:onerror", a1, a2)
        refreshIntervalId = setTimeout(refresh, 10)
    }
    xhr.ontimeout = function(a1, a2) {
        console.log("xhr:ontimeout", a1, a2)
        refreshIntervalId = setTimeout(refresh, 10)
    }
    xhr.timeout = 1000 * dT * 1.2
    xhr.send();
    return xhr;
}

var lineQueue = []

function updateBody(state) {
    lineQueue.push({
        element: makeLine(state.speed, state.angle),
        time: state.time
    })

    var arrow = document.getElementById("arrow")
    arrow.style.transform = "rotate(" + state.angle + "deg)"

    var speedIndicator = document.getElementById("speedIndicator")
    speedIndicator.innerHTML = "" + state.speed + " m/s"

    var directionTextIndicator = document.getElementById("directionTextIndicator")
    directionTextIndicator.innerHTML = "Fra " + getDirectionText(state.angle)
}

var angleBias = 0
var speedBias = 0

function makeLine(speed, angle) {
    if (this.prev === undefined) {
        this.prev = { x: 0, y: 0 }
    }

    // Rotate the chart so that North is up
    angle -= 360 / 4

    next = {
        x: speed * Math.cos((2 * Math.PI) / 360 * angle),
        y: speed * Math.sin((2 * Math.PI) / 360 * angle)
    }

    var line = document.createElementNS("http://www.w3.org/2000/svg", "line")
    line.setAttribute("x1", this.prev.x)
    line.setAttribute("y1", this.prev.y)
    line.setAttribute("x2", next.x)
    line.setAttribute("y2", next.y)

    this.prev.x = next.x
    this.prev.y = next.y

    return line
}
var svgGraph = null
var graphTimeIndicator = null

function drawNextLine(doStyleLines) {
    if (lineQueue.length > 0) {
        if (svgGraph === null || graphTimeIndicator === null) {
            svgGraph = document.getElementById("fancyWindGraph")
            graphTimeIndicator = document.getElementById("graphTimeIndicator")
        }

        if (doStyleLines !== false) {
            styleLines()
        }

        var line = lineQueue.shift()

        svgGraph.append(line.element)

        if (svgGraph.children.length > MAX_LINES) {
            svgGraph.children[0].remove()
        }

        line.time = line.time
            .replace("T", " ")
            .replace("\"", "")
            .replace(/\..*$/, "")

        graphTimeIndicator.innerHTML = line.time
    }
}

function styleLines() {
    var fancyWindGraph = document.getElementById("fancyWindGraph")
    var lines = fancyWindGraph.getElementsByTagName("line")
    // i=1 to skip first line that starts center
    for (var i = 1; i < lines.length; i++) {
        if (i < MAX_LINES) {
            var progress = (lines.length - i) / lines.length
            var colorFactor = 1 - theFormula(progress)
            var r = 128 + (96 * (colorFactor))
            var g = 128 + (96 * (colorFactor))
            var b = 256 - (96 * (colorFactor))
            var a = 1 - Math.pow(progress, 0.05) + 0.3
            lines[i].setAttribute("stroke", "rgba(" + r + "," + g + "," + b + ", " + a + ")")
        } else {
            lines[i].remove()
        }
    }
}

function theFormula(x) {
    return Math.pow(Math.log10(x), 4)
    //return Math.pow(x, 0.3)
}

var refreshIntervalId = setTimeout(refresh, 1000)
var drawNextLineIntervalId = setInterval(drawNextLine, 25)

function getDirectionText(angle) {
    if (angle < 0 + (22.5 / 2))
        return "nord"
    else if (angle < 22.5 + (22.5 / 2))
        return "nordnordøst"
    else if (angle < 45 + (22.5 / 2))
        return "nordøst"
    else if (angle < 67.5 + (22.5 / 2))
        return "østnordøst"
    else if (angle < 90 + (22.5 / 2))
        return "øst"
    else if (angle < 112.5 + (22.5 / 2))
        return "østsørøst"
    else if (angle < 135 + (22.5 / 2))
        return "sørøst"
    else if (angle < 157.5 + (22.5 / 2))
        return "sørsørøst"
    else if (angle < 180 + (22.5 / 2))
        return "sør"
    else if (angle < 202.5 + (22.5 / 2))
        return "sørsørvest"
    else if (angle < 225 + (22.5 / 2))
        return "sørvest"
    else if (angle < 247.5 + (22.5 / 2))
        return "vestsørvest"
    else if (angle < 270 + (22.5 / 2))
        return "vest"
    else if (angle < 292.5 + (22.5 / 2))
        return "vestnordvest"
    else if (angle < 315 + (22.5 / 2))
        return "nordvest"
    else if (angle < 337.5 + (22.5 / 2))
        return "nordnordvest"
    else
        return "nord"
}

/* For testing of the fancy graph.
// Copy into JS console in browser dev tools:
for(let i = 0; i < MAX_LINES; i++){
    updateBody({
        speed: (i / MAX_LINES)*25,
        angle: (i/MAX_LINES * 2) * 360,
        time: "testing"
    })
}
for(let i = MAX_LINES; i >= 0; i--){
    updateBody({
        speed: (i / MAX_LINES)*25,
        angle: 360 - (i/MAX_LINES * 2) * 360,
        time: "testing"
    })
}
*/