import React, { useRef, useEffect, useCallback } from 'react'
import { cn } from '@/lib/utils'

interface OTPInputProps {
  length?: number
  value: string
  onChange: (value: string) => void
  disabled?: boolean
  error?: string
  label?: string
}

export function OTPInput({
  length = 6,
  value,
  onChange,
  disabled = false,
  error,
  label = 'Código de verificação',
}: OTPInputProps) {
  const inputRefs = useRef<(HTMLInputElement | null)[]>([])
  const errorId = 'otp-error'

  useEffect(() => {
    inputRefs.current[0]?.focus()
  }, [])

  const handleChange = useCallback(
    (index: number, char: string) => {
      if (!/^\d*$/.test(char)) return

      const newValue = value.split('')
      newValue[index] = char
      const result = newValue.join('').slice(0, length)
      onChange(result)

      if (char && index < length - 1) {
        inputRefs.current[index + 1]?.focus()
      }
    },
    [value, onChange, length]
  )

  const handleKeyDown = useCallback(
    (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Backspace' && !value[index] && index > 0) {
        inputRefs.current[index - 1]?.focus()
        const newValue = value.split('')
        newValue[index - 1] = ''
        onChange(newValue.join(''))
      }
    },
    [value, onChange]
  )

  const handlePaste = useCallback(
    (e: React.ClipboardEvent) => {
      e.preventDefault()
      const pasted = e.clipboardData.getData('text').replace(/\D/g, '').slice(0, length)
      onChange(pasted)

      const focusIndex = Math.min(pasted.length, length - 1)
      inputRefs.current[focusIndex]?.focus()
    },
    [onChange, length]
  )

  return (
    <div>
      {/* Visually hidden group label for screen readers */}
      <p id="otp-group-label" className="sr-only">
        {label}
      </p>

      <div
        role="group"
        aria-labelledby="otp-group-label"
        aria-describedby={error ? errorId : undefined}
        className="flex gap-2 justify-center"
        onPaste={handlePaste}
      >
        {Array.from({ length }, (_, i) => (
          <input
            key={i}
            ref={(el) => {
              inputRefs.current[i] = el
            }}
            type="text"
            inputMode="numeric"
            maxLength={1}
            value={value[i] || ''}
            onChange={(e) => handleChange(i, e.target.value)}
            onKeyDown={(e) => handleKeyDown(i, e)}
            disabled={disabled}
            aria-label={`Dígito ${i + 1} de ${length}`}
            aria-invalid={!!error}
            aria-describedby={error ? errorId : undefined}
            autoComplete={i === 0 ? 'one-time-code' : 'off'}
            className={cn(
              'w-12 h-14 text-center text-2xl font-bold rounded-lg border-2 transition-all',
              'bg-white dark:bg-gray-800',
              'text-gray-900 dark:text-gray-100',
              'focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500',
              'disabled:opacity-50 disabled:cursor-not-allowed',
              error
                ? 'border-danger-500 focus:ring-danger-500'
                : 'border-gray-300 dark:border-gray-600'
            )}
          />
        ))}
      </div>

      {error && (
        <p
          id={errorId}
          role="alert"
          className="mt-2 text-sm text-center text-danger-600 dark:text-danger-400"
        >
          {error}
        </p>
      )}
    </div>
  )
}

export default OTPInput
