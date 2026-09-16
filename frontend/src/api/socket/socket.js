import { activePage } from '@/data/chatState';
import { addNotification } from '@/data/notifications';
import { handleIncomingNotification } from '@/data/notificationCount';

let ws = null;
const notificationDebounce = new Map();

export function connectToWS() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        return ws;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    
    ws = new WebSocket(`${protocol}//${window.location.host}/api/ws`);

    ws.onopen = () => {
        console.log('websocket connected');
    };

    ws.onmessage = (event) => {
        const payload = JSON.parse(event.data);
        switch (payload.type) {
            case 'notification':
                console.log('Notification:', payload.data);
                handleIncomingNotification(payload.data);
                break;

            case 'message': {
                const message = payload.data;

                console.log(message);
                console.log(activePage.value);
                console.log(activePage.value === 'chat:' + message.Sender.ID);

                if (activePage.value === 'chat:' + message.Sender.ID) {
                    window.dispatchEvent(
                        new CustomEvent('chat-message', {
                            detail: message
                        })
                    );
                } else {
                    const groupID = message.GroupID;
                    const now = Date.now();
                    const lastNotification = notificationDebounce.get(groupID);

                    if (
                        !lastNotification ||
                        now - lastNotification >= 30 * 60 * 1000
                    ) {
                        addNotification(
                            `New message from ${message.Sender.firstName || 'user'}`,
                            'message'
                        );

                        notificationDebounce.set(groupID, now);
                    }
                }

                break;
            }
        }
    };

    ws.onclose = () => {
        console.log('websocket disconnected');
        ws = null;
    };

    ws.onerror = (error) => {
        console.error('websocket error:', error);
    };

    return ws;
}

export function sendWS(payload) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        return;
    }

    ws.send(JSON.stringify(payload));
}