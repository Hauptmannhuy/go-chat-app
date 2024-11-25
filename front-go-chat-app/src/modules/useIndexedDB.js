import { useEffect, useRef, useState } from "react";

export function useIndexedDB() {
  const indexDB = useRef(null)

  const [cacheStatus, setCacheStatus] = useState(true)
  useEffect(() => {

   async function connectDB(){
     const openRequest = window.indexedDB.open('565763124')
     openRequest.addEventListener("error", () => {
       console.log("error connect")
      })
      
      openRequest.addEventListener("success", () => {

        indexDB.current = openRequest.result
        console.log("connected to IndexedDB")
      })

      
      openRequest.addEventListener("upgradeneeded", () => {
        const db = openRequest.result
        console.log("upgrading")

        const messageStore = db.createObjectStore("messages",{keyPath: "id", autoIncrement: true}) 
        messageStore.createIndex("message_id", "message_id", {unique: false})
        messageStore.createIndex("body", "body", {unique: false})
        messageStore.createIndex("chat_name", "chat_name", {unique: false})
        messageStore.createIndex("user_id", "user_id", {unique: false} )
        const privateChatStore = db.createObjectStore("private_chats",{keyPath: "id"}) 
        privateChatStore.createIndex("id", "id", {unique: true})
        privateChatStore.createIndex("user1_name", "user1_name", {unique: false})
        privateChatStore.createIndex("user2_name", "user2_name", {unique: false})

        const groupChatStore = db.createObjectStore("group_chats",{keyPath: "id"}) 
        groupChatStore.createIndex("id", "id", {unique: true})
        groupChatStore.createIndex("name", "name", {unique: false})
        groupChatStore.createIndex("creator_id", "creator_id", {unique:false})
        
        setCacheStatus(false)
        
      })

    }

    connectDB()
  }, [])

      /**
       * Returns objectStore.
       * @param {string} dbname
       * @returns {IDBObjectStore}
       */

   const initDBtransaction = (dbname) => { 
    const transaction = indexDB.current.transaction([dbname], "readwrite")
    const objectStore = transaction.objectStore(dbname)
    return objectStore
  }

  function cacheMessages(data){
    console.log("db:", indexDB)
    const messageStore = initDBtransaction("messages")
    console.log(data)
    // const parsedData = JSON.parse(data)
    const objectChatNames = Object.keys(data)
    for (let i = 0; i < objectChatNames.length; i++) {
      const chatName = objectChatNames[i]
      const chatObj = data[chatName]
      for (let j = 0; j < chatObj.length; j ++) {
        console.log(chatObj[j])
        const newMessageObject = {message_id: chatObj[j].message_id, body: chatObj[j].body, chat_name: chatObj[j].chat_name, username: chatObj[j].username }
        const addRequest = messageStore.add(newMessageObject)

        addRequest.addEventListener("error", () => {
          console.log("Error adding new item to DB")
          console.log(addRequest.error)
        })
      }
    }
  }

  function cachePrivateChats(data){
    
  }

  function cacheGroupChats(data){

  }

  return {cacheStatus, cacheMessages}
}