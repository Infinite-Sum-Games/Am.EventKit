import pandas as pd
import matplotlib.pyplot as plt

def analyze_usage(file_path="resource_usage_1.csv"):
    """
    Analyzes the resource usage data from the CSV file and generates plots.
    """
    try:
        # Read the CSV file
        df = pd.read_csv(file_path, encoding='utf-16')

        # Convert Timestamp to datetime objects
        df['Timestamp'] = pd.to_datetime(df['Timestamp'])

        # --- CPU Usage Analysis ---
        # Calculate the difference in CPU seconds between each interval
        df['CPU_Seconds_Delta'] = df['CPU_TotalSeconds'].diff()

        # Calculate the time delta in seconds between each timestamp
        df['Time_Delta_Seconds'] = df['Timestamp'].diff().dt.total_seconds()

        # Calculate CPU usage percentage
        # This is an approximation. For a more accurate measure, one would need to divide by the number of cores.
        # This calculation gives a measure of "CPU seconds per second".
        df['CPU_Usage'] = (df['CPU_Seconds_Delta'] / df['Time_Delta_Seconds']) * 100
        
        # Handle the first row which will be NaN after diff()
        df = df.dropna()

        # --- Generate Plots ---

        # 1. Memory Usage Plot
        plt.figure(figsize=(12, 6))
        plt.plot(df['Timestamp'], df['Memory_MB'], label='Memory Usage (MB)')
        plt.title('Memory Usage Over Time')
        plt.xlabel('Time')
        plt.ylabel('Memory (MB)')
        plt.grid(True)
        plt.legend()
        plt.xticks(rotation=45)
        plt.tight_layout()
        plt.savefig('memory_usage_1.png')
        print("Saved memory usage plot to memory_usage.png")

        # 2. CPU Usage Plot
        plt.figure(figsize=(12, 6))
        plt.plot(df['Timestamp'], df['CPU_Usage'], label='CPU Usage (%)', color='orange')
        plt.title('CPU Usage Over Time')
        plt.xlabel('Time')
        plt.ylabel('CPU Usage (%)')
        plt.grid(True)
        plt.legend()
        plt.xticks(rotation=45)
        plt.tight_layout()
        plt.savefig('cpu_usage_1.png')
        print("Saved CPU usage plot to cpu_usage.png")

        # --- Summary Statistics ---
        print("\n--- Resource Usage Report ---")
        print(f"Analysis of: {file_path}")
        print("\nMemory Usage (MB):")
        print(df['Memory_MB'].describe())
        print("\nCPU Usage (%):")
        print(df['CPU_Usage'].describe())
        print("\n-----------------------------")


    except FileNotFoundError:
        print(f"Error: The file '{file_path}' was not found.")
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    analyze_usage()
