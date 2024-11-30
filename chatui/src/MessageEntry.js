import React, { useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

var WS_URL = "ws://127.0.0.1:8080/ws";

const MessageEntry = ({ selectedChannel, onNewMessage }) => {
    const [text, setText] = useState('');

    if (selectedChannel) {
        // console.log('CHANNEL SELECTED => ', selectedChannel)
        WS_URL = `ws://127.0.0.1:8080/ws?roomID=${selectedChannel.id}`;
    }

    const { sendMessage, lastMessage, readyState } = useWebSocket(WS_URL, {
        onOpen: () => {
            console.log("WS connection made in MessageEntry")
        },
        share: true
    });

    const handleSendMessage = async () => {
        const userID = localStorage.getItem('userId');
        const userName = localStorage.getItem('userName');

        const response = await fetch('/messages', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                "channel_id": parseInt(selectedChannel.id),
                "user_id": parseInt(userID),
                text
            }),
        });

        if (response.ok) {
            const message = await response.json();
            onNewMessage({
                id: message.id,
                channel_id: selectedChannel,
                user_id: userID,
                user_name: userName,
                text
            });

            setText('');

            sendMessage(JSON.stringify({
                id: message.id,
                channel_id: selectedChannel,
                user_id: userID,
                user_name: userName,
                text
            }));

        } else {
            alert('Failed to send message');
        }
    };

    const handleKeyDown = (event) => {
        if (event.key === 'Enter' && !event.shiftKey) {
            handleSendMessage();
            event.preventDefault(); // Prevent the default behavior (newline)
        }
    };

    return (
        <div className="p-4 border-t flex">
            <input
                type="text"
                placeholder="Type a message..."
                value={text}
                onChange={(e) => setText(e.target.value)}
                onKeyDown={handleKeyDown}
                className="p-2 flex-grow border rounded-md mr-2"
            />
            <button
                onClick={handleSendMessage}
                className="p-2 bg-blue-500 text-white rounded-md"
            >
                Send
            </button>
        </div>
    );
};

export default MessageEntry;