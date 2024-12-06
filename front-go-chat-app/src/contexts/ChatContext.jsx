import { createContext } from "react";

export const ChatContext = createContext({
  chatManager:{
    chats: {},

  },
  messages: {},
  searchResults: {},
  searchProfileResults: {},

})

