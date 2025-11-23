// API module for backend communication

// Cache config at module level to prevent multiple API calls
let configPromise = null;

/**
 * Fetches public configuration from the server
 * @returns {Promise<{isDefaultPassword: boolean, fileExpiryMinutes: number, maxFileSizeMB: number}>}
 */
export async function getPublicConfig() {
	if (!configPromise) {
		configPromise = (async () => {
			try {
				const response = await fetch('/config');
				if (response.ok) {
					return await response.json();
				}
				return { isDefaultPassword: true, fileExpiryMinutes: 10, maxFileSizeMB: 100 };
			} catch (error) {
				console.error('Failed to fetch config:', error);
				return { isDefaultPassword: true, fileExpiryMinutes: 10, maxFileSizeMB: 100 };
			}
		})();
	}
	return configPromise;
}

/**
 * Validates the upload password with the server
 * @param {string} password - The password to validate
 * @returns {Promise<boolean>} - True if password is valid
 */
export async function validatePassword(password) {
	try {
		const response = await fetch('/login', {
			method: 'POST',
			headers: {
				'X-Upload-Password': password
			}
		});
		return response.ok;
	} catch (error) {
		console.error('Password validation error:', error);
		return false;
	}
}

/**
 * Uploads a file to the server
 * @param {File} file - The file to upload
 * @param {string} password - The upload password
 * @param {Function} onProgress - Progress callback (percentage)
 * @returns {Promise<{success: boolean, data?: object, error?: string}>}
 */
export function uploadFile(file, password, onProgress) {
	return new Promise((resolve) => {
		const formData = new FormData();
		formData.append('file', file);

		const xhr = new XMLHttpRequest();

		// Track upload progress
		xhr.upload.addEventListener('progress', (e) => {
			if (e.lengthComputable && onProgress) {
				const percentage = Math.round((e.loaded / e.total) * 100);
				onProgress(percentage);
			}
		});

		// Handle completion
		xhr.addEventListener('load', () => {
			if (xhr.status === 200) {
				const response = xhr.responseText;
				const urlMatch = response.match(/Download URL: (http[^\s]+)/);
				const curlMatch = response.match(/cURL command: (.+)/);
				const nameMatch = response.match(/Original name: ([^\n]+)/);
				const sizeMatch = response.match(/File size: ([^\n]+)/);

				if (urlMatch) {
					const downloadURL = urlMatch[1];
					const curlCommand = curlMatch ? curlMatch[1] : '';
					const fileName = nameMatch ? nameMatch[1] : '';
					const fileSize = sizeMatch ? sizeMatch[1] : '';
					const fileIDMatch = downloadURL.match(/\/download\/([^\/\s]+)/);
					const fileID = fileIDMatch ? fileIDMatch[1] : null;

					resolve({
						success: true,
						data: {
							fileName,
							fileSize,
							downloadURL,
							curlCommand,
							fileID
						}
					});
				} else {
					resolve({
						success: false,
						error: 'Upload successful but couldn\'t parse response'
					});
				}
			} else if (xhr.status === 401) {
				resolve({
					success: false,
					error: 'Unauthorized: Invalid password',
					unauthorized: true
				});
			} else {
				// Try to get the error message from the response body
				let errorMessage = xhr.responseText.trim() || xhr.statusText;
				// If response is too long, truncate it (likely HTML error page)
				if (errorMessage.length > 200 || errorMessage.includes('<html')) {
					errorMessage = xhr.statusText;
				}
				resolve({
					success: false,
					error: `Upload failed: ${errorMessage}`
				});
			}
		});

		// Handle errors
		xhr.addEventListener('error', () => {
			resolve({
				success: false,
				error: 'Upload failed: Network error'
			});
		});

		// Send request
		xhr.open('POST', '/');
		xhr.setRequestHeader('X-Upload-Password', password);
		xhr.send(formData);
	});
}
