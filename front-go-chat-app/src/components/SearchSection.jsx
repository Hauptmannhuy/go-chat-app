import { useContext, useRef, useState } from "react"
import { SearchContext } from "../contexts/SearchContext"

export default function SearchSection({searchHandler}){
  const {searchResults} = useContext(SearchContext)
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