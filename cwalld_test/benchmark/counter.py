import csv
import matplotlib.pyplot as plt
import numpy as np

with open("results.csv", "r") as f:
    reader = csv.reader(f)
    data = [int(x) for x in list(reader)[0]]

occurences = {}

for i in range(0, 501):
    data.append(i)

data = sorted(data)

for i in range(0, len(data)):
    if data[i] not in occurences:
        occurences[data[i]] = -1
    occurences[data[i]] += 1

minimum = 500
maximum = 0

for i in range(0, len(data)):
    if occurences[data[i]] != 15:
        if data[i] < minimum:
            minimum = data[i] 
            amount = occurences[data[i]]

    if data[i] > maximum and occurences[data[i]] != 0:
        maximum = data[i]

    if occurences[data[i]] >= 16:
        print(data[i])

print(f"Quickest enforcal found at {minimum} with {amount} / 15 getting through")
print(f"Longest failure to enforce found at {maximum}")

threshold = min(t for t, f in occurences.items() if f == 0)

sleep_durations = list(occurences.keys())

for i in range(0, len(sleep_durations)):
    sleep_durations[i] = sleep_durations[i] * 10

sleep_past_thresh = [x for x in sleep_durations if x > (threshold * 10)]

for i in range((threshold), 500):
    if i > threshold and i*10 not in sleep_past_thresh:
        for j in range(0,9):
            sleep_past_thresh.append(0)

    elif i > threshold:
        for j in range(0, 8):
            sleep_past_thresh.append(0)


sleep_past_thresh.pop() # i know theres 3 extra ones that happened past the threshold that arent one
sleep_past_thresh.pop()
sleep_past_thresh.pop()

non_zero = np.nonzero(sleep_past_thresh)
percent = (len(non_zero) / len(sleep_past_thresh)) * 100
print(f"Past a threshold of {threshold} {percent}% of breaches dont get enforced")

# mean = np.mean(sleep_past_thresh)
# stdev = np.std(sleep_past_thresh)
# median = np.median(sleep_past_thresh)
# print(f"Past a threshold of {threshold} mean of {mean} stdev of {stdev} median of {median}") # honestly these arent very important

occurences_list = list(occurences.values())

# print(sleep_durations)
# print(occurences)
# print(occurences_list)

plt.plot(sleep_durations, occurences_list, marker = 'o', linewidth=1, markersize=2)
plt.xlabel("Time Slept (microseconds)")
plt.ylabel("Failures to enforce")
plt.ylim(0, 16)
plt.xlim(0, 5000)
plt.tight_layout()
plt.grid(True, linestyle="--", alpha=0.5)
plt.savefig("graph.png", dpi=150)

print("Graph created!")
