'use client';

import { useEffect, useState } from "react";
import { v4 } from "uuid"
import styles from './page.module.css'; 
import { CalculationResult, HistoryData, ComputeResponse, HistoryResponse } from "./model";
import { InfoCard } from "./components/info-card";
import { HistoryCard } from "./components/history-card";
import { fetchWithErrorHandle, sessionIdIsValid } from "./util";
import { Result } from "./components/result";
import { SessionIdCard } from "./components/session-id-card";
import { InputCard } from "./components/input-card";

const SESSION_ID_KEY = "session-id";


const API_BASE = "/api";


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

              <InputCard currentResult={currentResult} onCalculate={calculate}></InputCard>

              { sessionId ?        
                 <>
                 <SessionIdCard onChangeSessionId={changeSessionId} sessionId={sessionId}></SessionIdCard> 
                 <HistoryCard historyData={historyData} onFetchNextHistoryPage={fetchNextHistoryPage} sessionId={sessionId}></HistoryCard>
                 </>
                 : null}
              <InfoCard></InfoCard>
          </main>  
    </>
  );
}
