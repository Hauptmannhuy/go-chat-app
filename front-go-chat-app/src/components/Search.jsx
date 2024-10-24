import { useRef, useState } from "react"

export default function Search({searchHandler}){
  const timeout = useRef(null)
 
  function typeEventHandler(input){
    if (timeout.current){
      clearTimeout(timeout.current)
    }
    timeout.current = setTimeout(()=>{
      searchHandler(input)
    }, 400)
  }
  return (<>
  <input type="text" className="searchInput" onChange={(e) => {typeEventHandler(e.target.value)}} />
  </>)
}