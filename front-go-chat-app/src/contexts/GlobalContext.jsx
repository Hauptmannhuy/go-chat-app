import { createContext } from "react";

export const GlobalContext = createContext({
    socket: null, 
    sendMessage: null,
    selectChat: null,
    searchResults: {},
    searchProfileResults: {},
    selectedChat: {chat_name: null, id: null, participation: null, type: ''},
    chats: {},
    messages: {}
})