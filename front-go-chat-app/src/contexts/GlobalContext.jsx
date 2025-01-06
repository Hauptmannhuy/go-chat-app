import { createContext } from "react";

export const GlobalContext = createContext({
    socket: null, 
    writeToSocket: (message) => {},
    selectChat: null,
    searchResults: {},
    searchProfileResults: {},
    selectedChat: {name: null, id: null, participation: null, type: ''},
    chats: {},
    messages: {},
    onlineStatus: (name) => {},
})