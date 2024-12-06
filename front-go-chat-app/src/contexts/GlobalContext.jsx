import { createContext } from "react";

export const GlobalContext = createContext({
    socket: null, 
    sendMessage: null,
    selectChat: null,
    selectedChat: {chat_name: null, chat_id: null},
    chats: {},
    messages: {}
})