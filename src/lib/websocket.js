// WebSocket module for real-time download notifications

/**
 * Creates a WebSocket connection for file download notifications
 * @param {string} fileID - The file ID to monitor
 * @param {Function} onDownloaded - Callback when file is downloaded
 * @returns {WebSocket} - The WebSocket connection
 */
export function connectToFileNotifications(fileID, onDownloaded) {
	const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	const wsUrl = `${protocol}//${window.location.host}/ws/${fileID}`;

	const websocket = new WebSocket(wsUrl);

	websocket.onopen = () => {
		console.log('WebSocket connected for file:', fileID);
	};

	websocket.onmessage = (event) => {
		try {
			const data = JSON.parse(event.data);
			if (data.downloaded && onDownloaded) {
				onDownloaded();
			}
		} catch (error) {
			console.error('WebSocket message error:', error);
		}
	};

	websocket.onerror = (error) => {
		console.error('WebSocket error:', error);
	};

	websocket.onclose = () => {
		console.log('WebSocket closed');
	};

	return websocket;
}

/**
 * Closes a WebSocket connection
 * @param {WebSocket} websocket - The WebSocket to close
 */
export function closeWebSocket(websocket) {
	if (websocket) {
		websocket.close();
	}
}
