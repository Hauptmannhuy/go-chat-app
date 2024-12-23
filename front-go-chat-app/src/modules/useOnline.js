import { useState } from "react"

export default function useOnline(){
  const [onlineUsers, setOnlineUsers] = useState({})

  function changeOnlineStatus(newStatusMessage) {
    console.log(newStatusMessage)
    setOnlineUsers((oldStatus) => {
    const newOnlineStatus = {...oldStatus}
    const keys = Object.keys(newStatusMessage)
    for (const key of keys) {
      newOnlineStatus[key] = newStatusMessage[key]
    }
    console.log(onlineUsers)
    return newOnlineStatus
    })  
  }
  
  const onlineStatus = (name) => (onlineUsers[name])

  return {onlineUsers, changeOnlineStatus, onlineStatus}
}