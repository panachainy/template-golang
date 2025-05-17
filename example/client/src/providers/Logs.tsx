import type React from 'react'
import { createContext, useContext, useState } from 'react'

interface LogsContextType {
  logs: string[]
  addLog: (log: string) => void
}

const LogsContext = createContext<LogsContextType | undefined>(undefined)

export const LogsProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [logs, setLogs] = useState<string[]>([])

  const addLog = (log: string) => {
    setLogs((prevLogs) => [...prevLogs, log])
  }

  return (
    <LogsContext.Provider value={{ logs, addLog }}>
      {children}
      <LogsComponent logs={logs} />
    </LogsContext.Provider>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export const useLogs = (): LogsContextType => {
  const context = useContext(LogsContext)
  if (!context) {
    throw new Error('useLogs must be used within a LogsProvider')
  }
  return context
}

const LogsComponent: React.FC<{ logs: string[] }> = ({ logs }) => {
  return (
    <div
      style={{
        position: 'absolute',
        top: 0,
        right: 0,
        backgroundColor: 'rgba(255, 255, 255, 0.8)',
        padding: '10px',
        border: '1px solid black',
        zIndex: 1000,
      }}
    >
      <strong>Logs:</strong>
      <ul>
        {logs.map((log) => (
          <li key={`${log}`}>{log}</li>
        ))}
      </ul>
    </div>
  )
}
