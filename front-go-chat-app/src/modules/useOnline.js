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
    return newOnlineStatus
    })  
  }
  
  const onlineStatus = (name) => {
    console.log(onlineUsers[name])
    return onlineUsers[name]}

  return {onlineUsers, changeOnlineStatus, onlineStatus}
}