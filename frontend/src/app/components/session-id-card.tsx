import { useState } from "react";
import styles from "../page.module.css"
import { sessionIdIsValid } from "../util";


export function SessionIdCard({onChangeSessionId, sessionId}: 
    {onChangeSessionId: (newSessionId: string) => void, sessionId: string}) {
    const [sessionIdInputValue, setSessionIdInputValue] = useState<string>("");
  
    function onKeyUp(e: any) {
      if (e && e.key == 'Enter') {
        if (sessionIdIsValid(sessionIdInputValue)) {
            onChangeSessionId(sessionIdInputValue);
        } else {
            alert("Session id is not valid!");
        }
      }
    }
  
    return (
      <div className={styles.gridCard}>
      <h2>Change session id</h2>
      <div>
        <input value={sessionIdInputValue} onChange={ e => setSessionIdInputValue(e.currentTarget.value)} onKeyUp={onKeyUp}></input>
        <button onClick={() => onChangeSessionId(sessionIdInputValue)}>Change</button>
      </div>
      </div>
    )
  }