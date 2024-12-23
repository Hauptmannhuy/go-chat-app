import { createContext } from "react";

export const GlobalContext = createContext({
    socket: null, 
    sendMessage: null,
    selectChat: null,
    searchResults: {},
    searchProfileResults: {},
    selectedChat: {name: null, id: null, participation: null, type: ''},
    chats: {},
    messages: {},
    onlineStatus: (name) => {},
})