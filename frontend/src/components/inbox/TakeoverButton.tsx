import Button from '../ui/button/Button';
import { Power } from 'lucide-react';

interface TakeoverButtonProps {
  humanMode: boolean;
  onClick: () => void;
  className?: string;
}

export default function TakeoverButton({ humanMode, onClick, className }: TakeoverButtonProps) {
  return (
    <div data-testid="takeover-button" className={className}>
      <Button 
        variant={humanMode ? "primary" : "outline"}
        onClick={onClick}
        startIcon={<Power size={16} />}
      >
        {humanMode ? "Human Mode Active" : "Take Over (Human Mode)"}
      </Button>
    </div>
  );
}