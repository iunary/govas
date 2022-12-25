const canvas = document.querySelector("#canvas");
const ctx = canvas.getContext("2d");
let isDrawing = false;

// canvas
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

ctx.lineJoin = "round";
ctx.lineCap = "round";
ctx.lineWidth = 4;
ctx.strokeStyle = "#000000";

// websocket
var socket = new WebSocket(`ws://${document.location.host}/ws`);
socket.onmessage = (event) => {
	try {
		const data = JSON.parse(event.data);
		console.log("data ::: ", data);
		switch (data.action) {
			case "draw":
				isDrawing = data.isDrawing;
				draw(data.x, data.y, data.type);
				break;
			case "clear":
				ctx.clearRect(0, 0, canvas.width, canvas.height);
				break;
		}
	} catch (error) {
		console.log(error);
	}
};

socket.onerror = (event) => {
	console.log("error", event);
};

function draw(x, y, type) {
	if (!isDrawing) return;
	if (type === "mousedown") {
		ctx.beginPath();
		return ctx.moveTo(x, y);
	} else if (type === "mousemove") {
		ctx.lineTo(x, y);
		return ctx.stroke();
	} else {
		return ctx.closePath();
	}
}

canvas.addEventListener("mousedown", (e) => {
	isDrawing = true;
	socket.send(
		JSON.stringify({
			action: "draw",
			x: e.offsetX,
			y: e.offsetY,
			type: e.type,
			isDrawing: isDrawing,
		}),
	);
});

canvas.addEventListener("mousemove", (e) => {
	if (!isDrawing) return;
	socket.send(
		JSON.stringify({
			action: "draw",
			x: e.offsetX,
			y: e.offsetY,
			type: e.type,
			isDrawing: isDrawing,
		}),
	);
});

canvas.addEventListener("mouseup", () => (isDrawing = false));
canvas.addEventListener("mouseout", () => (isDrawing = false));

// toolbox
let clearbtn = document.getElementById("clearbtn");
clearbtn.addEventListener("click", (e) => {
	socket.send(
		JSON.stringify({
			action: "clear",
		}),
	);
});
