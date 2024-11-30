import React, { useState, useEffect, useRef } from 'react';
import MessageEntry from './MessageEntry';
import useWebSocket, { ReadyState } from 'react-use-websocket';

var WS_URL = "ws://127.0.0.1:8080/ws";

//  Add websocket connection here ???
const MessagesPanel = ({ selectedChannel }) => {
    const [messages, setMessages] = useState([]);
    const lastMessageIdRef = useRef(null); // Keep track of the last message ID

    if (selectedChannel) {
        // console.log('CHANNEL SELECTED => ', selectedChannel)
        WS_URL = `ws://127.0.0.1:8080/ws?roomID=${selectedChannel.id}`;
    }

    const { sendMessage, lastMessage, readyState } = useWebSocket(WS_URL, {
        onOpen: () => {
            console.log("WS connection made in MessagesPanel")
        },
        share: true
    });

    useEffect(() => {
        if (!selectedChannel) return;

        let isMounted = true; // flag to prevent state updates after unmount

        const fetchMessages = async () => {
            const response = await fetch(`/messages?channelID=${selectedChannel.id}`);
            const data = await response.json();
            if (isMounted) {
                let messageData = data || [];
                setMessages(messageData);
                lastMessageIdRef.current = messageData.length > 0 ? messageData[messageData.length - 1].id : null;
            }
        };

        fetchMessages();

        if (lastMessage !== null) {
            console.log("LAST MESSAGE: ", lastMessage);
            if (isMounted) {
                let newMessage = JSON.parse(lastMessage.data);

                setMessages((prev) => {
                    prev.concat(newMessage);
                    lastMessageIdRef.current = newMessage.id;
                    return messages;
                });
            }
        }

        return () => {
            isMounted = false; // prevent further state updates
        };
    }, [selectedChannel, lastMessage]);

    return (
        <div className="flex flex-col h-full">
            {selectedChannel && (
                <div className="bg-gray-700 text-white p-2">
                    Messages for {selectedChannel.name}
                </div>
            )}
            <div className={`overflow-auto flex-grow ${selectedChannel && messages.length === 0 ? 'flex items-center justify-center' : ''}`}>
                {selectedChannel ? (
                    messages.length > 0 ? (
                        messages.map((message) => (
                            <div key={message.id} className="p-2 border-b">
                                <strong>{message.user_name}</strong>: {message.text}
                            </div>
                        ))
                    ) : (
                        <div className="text-center text-gray-600">
                            No messages yet! Why not send one?
                        </div>
                    )
                ) : (
                    <div className="p-2">Please select or create a channel</div>
                )}
            </div>
            {selectedChannel && (
                <MessageEntry
                    selectedChannel={selectedChannel}
                    onNewMessage={(message) => {
                        lastMessageIdRef.current = message.id;
                        setMessages([...messages, message])
                    }
                    }
                />
            )}
        </div>
    );
};

export default MessagesPanel;