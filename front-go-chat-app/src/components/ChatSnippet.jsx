import { useContext } from "react"
import { GlobalContext } from "../contexts/GlobalContext"
import { nameFormatter } from "../modules/nameFormatter"

export function ChatSnippet({name, onSelect, currentChats}) {
  const {messages} = useContext(GlobalContext)
  
  const dialogue = messages[name]
  if (dialogue && dialogue.length > 0) {
    return (
      <div className="snippet-container" onClick={() => onSelect(currentChats[name])}>
        <div className="chatname-snippet"> {nameFormatter(name)}</div>
        <div className="chatmessage-snippet"> {dialogue[dialogue.length-1].body}</div>
      </div>
      )
    }  else {
      return (
        <div className="snippet-container" onClick={() => onSelect(currentChats[name])}>
          <div className="chatname-snippet">{nameFormatter(name)}</div>
        </div>
      )
    }
}