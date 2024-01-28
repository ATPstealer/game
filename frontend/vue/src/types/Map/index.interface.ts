export interface Cell {
  cellName: string;
  x: number;
  y: number;
  surfaceImagePath: string;
  square: number;
  pollution: number;
  population: number;
  education: number;
  crime: number;
  medicine: number;
  averageSalary: number;
}

export interface CellOwners {
  nickName: string;
  square: number;
  x: number;
  y: number;
}
