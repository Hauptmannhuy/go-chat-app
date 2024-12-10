import { useContext } from "react"
import { GlobalContext } from "../contexts/GlobalContext"
import { useAuth } from "../modules/useAuth"

export function ChatSnippet({name}) {
  const {getUsername} = useAuth()
  const {messages} = useContext(GlobalContext)
  const nameFormatter = (name) => (name.split('_').filter((el) => (el != getUsername())).join(' '))
  const dialogue = messages[name]
  // const lastMessage = dialogue[dialogue.length-1]
  if (dialogue && dialogue.length > 0) {
    return (
      <div className="snippet-container">
        <div className="chatname-snippet"> {nameFormatter(name)}</div>
        <div className="chatmessage-snippet"> {dialogue[dialogue.length-1].body}</div>
      </div>
      )
    }  else {
      return (
        <div className="snippet-container">
          <div className="chatname-snippet">{nameFormatter(name)}</div>
        </div>
      )
    }
}