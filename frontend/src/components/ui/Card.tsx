/**
 * Card — adapter wrapper around @iit/ui Card.
 * Maps local `hover`/`padding` props to the @iit/ui API.
 */
import React from 'react'
import { Card as IITCard, CardHeader, CardTitle, CardDescription, CardBody, CardFooter } from '@iit/ui'
import { cn } from '@/lib/utils'

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  hover?: boolean
  hoverable?: boolean
  padding?: 'none' | 'sm' | 'md' | 'lg'
}

const paddings = {
  none: 'p-0',
  sm: 'p-3',
  md: 'p-4',
  lg: 'p-6',
}

const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ children, hover, hoverable, padding = 'md', className, ...props }, ref) => {
    return (
      <IITCard
        ref={ref}
        hoverable={hoverable || hover}
        className={cn(
          'bg-white border-gray-200 text-gray-900',
          paddings[padding],
          className
        )}
        {...props}
      >
        {children}
      </IITCard>
    )
  }
)

Card.displayName = 'Card'

export { Card, CardHeader, CardTitle, CardDescription, CardBody, CardFooter }
export default Card
