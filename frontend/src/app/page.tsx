'use client';

import { KeyboardEventHandler, useEffect, useState } from "react";
import { v4 } from "uuid"
import styles from './page.module.css'; 
import { MathJax } from "better-react-mathjax";

const SESSION_ID_KEY = "session-id";

type CalculationResult = {
  input: string,
  result?: RationalNumber
  errorId?: string,
}

type RationalNumber =  {
  num: string,
  denom: string, 
  estimate: string
}

type HistoryRow = {
  calculationId: number,
  calculation: CalculationResult,
}

type HistoryData = {
  nextUrl?: string,
  items: HistoryRow[]
}

type ComputeResponse = {
  error?: string,
  historyRow: HistoryRow
}


const API_BASE = "/api";

type HistoryResponse = {
  self: string,
  first: string,
  prev?: string,
  next?: string,
  last: string,
  items: HistoryRow[]
  error: string,
}

function Result({calculationResult} : {calculationResult: CalculationResult}) {

  function  getMathJaxString() {

    if (!calculationResult.result) {
      return "";
    }

    if (calculationResult.result.denom === "1") {
      return `\\[${calculationResult.result.num}\\]`
    }
    
    return  `\\[\\frac{${calculationResult.result.num}}{${calculationResult.result.denom}}  \\approx ${calculationResult.result.estimate}\\]`;
  }
  return (
      <div className={styles.center}>
        {
          calculationResult.errorId ? <p>{calculationResult.errorId}</p> 
          : <MathJax inline dynamic>{getMathJaxString()}</MathJax> 
        }
        </div>
  )
}

function SessionIdTool({onChangeSessionId, sessionId}: 
  {onChangeSessionId: (newSessionId: string) => void, sessionId: string}) {
  const [sessionIdInputValue, setSessionIdInputValue] = useState<string>("");

  function onKeyUp(e: any) {
    if (e && e.key == 'Enter') {
      onChangeSessionId(sessionIdInputValue);
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

function InfoCard() {
  return (
    <div className={styles.gridCard}>
      <h2> Info </h2>
      <p>This is a calculator for computing expressions with rational numbers.</p>
      You can input
      <ul>
        <li> Integer constants, eg. 3, 5, 0, -0, -432</li>
        <li> The usual 4 operations, eg. 3+4, 4-2, 3*2, 5/2</li>
        <li> Any expressions of the above using brackets, eg. (4+3)*2/(5-2)</li>
      </ul>
      You <strong>can not</strong> input 
      <ul>
        <li> The unary operator '-' when not used in the context of a negative constant, eg. -(2+3), --5</li>
        <li> Any decimal representation of numbers, eg. 3.14, 100,00</li>
      </ul>

      <p>The calculator respects the order of operations an multiplication and division is evaluated left to right. For instance,</p>
      <p>3*4/2*5 = ((3*4)/2)*5</p>
      
    </div>
  )
}


function HistorySection({historyData, onFetchNextHistoryPage, sessionId}: 
  {historyData: HistoryData, onFetchNextHistoryPage: () => Promise<void>, sessionId: string}) {
  return (
    <div className={`${styles.historySection} ${styles.gridCard}`}>
      <h2>History</h2>
      <p>Session id: {sessionId}</p>
      <table className={styles.historyTable}>
            
        <thead>
        <tr>
          <th>Input</th>
          <th>Output</th>
        </tr>
        </thead>
        <tbody>
        {
          historyData.items.map( item => 
            <tr key={item.calculationId} className={styles.historyRow}>
              <td style={{textAlign: 'center'}}>{item.calculation.input}</td>

              <td style={{textAlign: 'center'}}>
                <Result calculationResult={item.calculation}></Result>
              </td>
            </tr>
          )
        }
        </tbody>
      </table>
      {
        historyData.nextUrl ? <button onClick={() => onFetchNextHistoryPage()}>Fetch more..</button> : null
      }
    </div>
  )
}

function InputSection({onCalculate, currentResult}: 
  {onCalculate: (input: string) => Promise<void>, currentResult?: CalculationResult | null}) {

  const [calculationInputValue, setCalculationInputValue] = useState<string>("");

  function onKeyUp(e: any) {
    if (e && e.key == 'Enter') {
      onCalculate(calculationInputValue)
    }
  }

  return (
    <>

  
     <div className={`${styles.inputSection} ${styles.gridCard}`}>
        <h2 >Input</h2>
        <div className={styles.inputMain}>
          <div className={styles.inputForm}>
            <input placeholder="1+(3+4)" 
              value={calculationInputValue} 
              onChange={ e => setCalculationInputValue(e.currentTarget.value)} 
              onKeyUp={(e) => onKeyUp(e)}></input>
            <button onClick={() => onCalculate(calculationInputValue)}>Calculate!</button>

          </div>  
          { currentResult ? <Result calculationResult={currentResult}></Result> : null}
        </div>
    </div>
    </>
  )
}


async function fetchWithErrorHandle<T>(
  input: string | URL | globalThis.Request,
  init?: RequestInit,
): Promise<T> {
  try {
    const response = await fetch(input, init)

    if (!response.ok) {
      alert(`Received response status: ${response.status} with body ${JSON.stringify(await response.json())}`)
      throw new Error(`Response not ok, status: ${response.status}`)
    }

    return await response.json()

  } catch(e) {
    alert("Received error from backend!");
    throw e;
  }
}

function sessionIdIsValid(sessionId: string): boolean {
  return /^[a-zA-Z0-9\-]*$/.test(sessionId) && sessionId.length > 1 && sessionId.length < 100
}

export default function Home() {
  const [sessionId, setSessionId] = useState<string|null>(null);
  const [historyData, setHistoryData] = useState<HistoryData>({items: []})
  const [currentResult, setCurrentResult] = useState<CalculationResult | null>();

  useEffect(() => {
    let sessionId = window.localStorage.getItem(SESSION_ID_KEY) ?? v4()

    if (!sessionIdIsValid(sessionId)) {
      sessionId = v4()
    }
    changeSessionId(sessionId)
  }, [])

  
  useEffect(() => {
    if (sessionId != null) {
      fetchHistoryForNewSessionId();
    }
  }, [sessionId])

  async function fetchHistoryForNewSessionId() {
    const historyResponse: HistoryResponse = await fetchWithErrorHandle(`/api/sessions/${sessionId}/history`)
    
    setHistoryData({
      nextUrl: historyResponse.next,
      items: historyResponse.items
    })
  }

  async function fetchNextHistoryPage() {
    if (historyData.nextUrl) {
      const historyResponse: HistoryResponse = await fetchWithErrorHandle(`${API_BASE}${historyData.nextUrl}`)
      
      setHistoryData((previousHistoryData) => {
        return {
          nextUrl: historyResponse.next,
          items: previousHistoryData.items.concat(historyResponse.items)
        }
      });
    }
  }

  async function changeSessionId(sessionId: string) {
    window.localStorage.setItem(SESSION_ID_KEY, sessionId);
    setSessionId(sessionId);
    setCurrentResult(null);
  }

  async function calculate(input: string) {

    const ComputeResponse: ComputeResponse = await fetchWithErrorHandle(`/api/sessions/${sessionId}/compute`, {
        method: 'POST',
        headers: { "Content-Type": "application/json"},
        body: JSON.stringify({
          "input": input
        })
      })

    setCurrentResult(ComputeResponse.historyRow.calculation)
    await fetchHistoryForNewSessionId()
  }

  

  return (
    <>
          <div className={styles.topBar}>

            <h1 className={styles.header}>Calculator</h1>
          </div>

          <main className={styles.main}>


              <InputSection currentResult={currentResult} onCalculate={calculate}></InputSection>

              { sessionId ?         <SessionIdTool onChangeSessionId={changeSessionId} sessionId={sessionId}></SessionIdTool> : null}
             

              {sessionId ? <HistorySection historyData={historyData} onFetchNextHistoryPage={fetchNextHistoryPage} sessionId={sessionId}></HistorySection>: null }
              


            <InfoCard></InfoCard>
          </main>
          
      
        

    </>
  );
}
