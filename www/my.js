function refresh() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/json");
    xhr.onload = updateBody;
    xhr.send();
    return xhr;
}

var lineQueue = []

function updateBody(response) {
    var data = JSON.parse(response.target.response)
    console.log(data)

    lineQueue.push(makeLine(data.speed, data.angle))

    var arrow = document.getElementById("arrow")
    arrow.style.transform = "rotate(" + data.angle + "deg)"

    var speedIndicator = document.getElementById("speedIndicator")
    speedIndicator.innerHTML = data.speed

    var timeIndicator = document.getElementById("timeIndicator")
    var split = data.updated.split(" ")
    timeIndicator.innerHTML = "" + split[0] + " " + split[1]
}

var angleBias = 0
var speedBias = 0

function makeLine(speed, angle) {
    if (this.prev === undefined) {
        this.prev = { x: 0, y: 0 }
    }

    speed = 25.0
    speedBias += (Math.random() * 2.5) - 1.25
    speedBias -= speedBias / 1000
    speed += speedBias
    var scale = 2 * speed

    angleBias += (Math.random() * 5) - 2.5
    angleBias -= angleBias / 1000
    angle += angleBias

    // Rotate the chart so that North is up
    angle -= 360 / 4

    next = {
        x: scale * Math.cos((2 * Math.PI) / 360 * angle),
        y: scale * Math.sin((2 * Math.PI) / 360 * angle)
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

function linePopper() {
    if (lineQueue.length > 0) {
        styleLines()

        var line = lineQueue.pop()

        var svgGraph = document.getElementById("fancyWindGraph")
        svgGraph.append(line)
    }
}

function styleLines() {
    var fancyWindGraph = document.getElementById("fancyWindGraph")
    var lines = fancyWindGraph.getElementsByTagName("line")
    // i=1 to skip first line that starts center
    for (var i = 1; i < lines.length; i++) {
        if (i < 1000) {
            var r = 128 + (128 * (lines.length - i)/lines.length)
            var g = 128 + (128 * (lines.length - i)/lines.length)
            var b = 256 - (128 * (lines.length - i)/lines.length)
            var a = theFormula((lines.length - i) / lines.length)
            lines[i].setAttribute("stroke", "rgba(" + r + "," + g + "," + b + ", " + a + ")")
        } else {
            lines[i].remove()
        }
    }
}

function theFormula(x) {
    return Math.pow((1 - Math.log10(x)) / 2, 4)
}

var refreshIntervalId = setInterval(refresh, 2000)
var linePopperIntervalId = setInterval(linePopper, 20)